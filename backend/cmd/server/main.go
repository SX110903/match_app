package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/config"
	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/email"
	"github.com/SX110903/match_app/backend/internal/handler"
	"github.com/SX110903/match_app/backend/internal/middleware"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

func main() {
	// 1. Load config (singleton - panics if required env vars are missing)
	cfg := config.Get()
	logger.Init(cfg.Server.Env)

	// 2. Connect to MySQL
	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to MySQL")
		os.Exit(1)
	}
	defer db.Close()
	logger.Info().Msg("connected to MySQL")

	// 3. Connect to Redis
	redisOpts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("invalid Redis URL")
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOpts)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to Redis")
		os.Exit(1)
	}
	defer redisClient.Close()
	logger.Info().Msg("connected to Redis")

	// 4. Instantiate auth services
	jwtSvc, err := auth.NewJWTService(cfg.JWT)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize JWT service")
		os.Exit(1)
	}
	totpSvc := auth.NewTOTPService()
	blacklist := auth.NewTokenBlacklist(redisClient)

	// 5. Instantiate repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	// 6. Instantiate services
	emailSvc := email.NewSMTPSender(cfg.Email)
	authSvc := service.NewAuthService(userRepo, tokenRepo, emailSvc, jwtSvc, totpSvc, blacklist, cfg)
	userSvc := service.NewUserService(userRepo, profileRepo)
	matchSvc := service.NewMatchService(matchRepo, profileRepo)

	// 7. Instantiate handlers
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	matchHandler := handler.NewMatchHandler(matchSvc)

	// 8. Register routes
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS(cfg.Security.AllowedOrigins))
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	// Global rate limit: 100 req/s per IP
	r.Use(middleware.NewIPRateLimiter(redisClient, 100, time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public with specific rate limits)
		r.Route("/auth", func(r chi.Router) {
			r.With(middleware.NewEndpointRateLimiter(redisClient, "register", 3, time.Hour)).
				Post("/register", authHandler.Register)

			r.With(middleware.NewEndpointRateLimiter(redisClient, "login", 5, 15*time.Minute)).
				Post("/login", authHandler.Login)

			r.With(middleware.NewEndpointRateLimiter(redisClient, "login_2fa", 10, 5*time.Minute)).
				Post("/login/2fa", authHandler.LoginWith2FA)

			r.Post("/logout", authHandler.Logout)
			r.Post("/refresh", authHandler.RefreshToken)
			r.Post("/verify-email", authHandler.VerifyEmail)

			r.With(middleware.NewEndpointRateLimiter(redisClient, "forgot_password", 3, time.Hour)).
				Post("/forgot-password", authHandler.ForgotPassword)

			r.Post("/reset-password", authHandler.ResetPassword)
		})

		// Protected routes
		authRequired := middleware.RequireAuth(jwtSvc, blacklist)

		r.Route("/users", func(r chi.Router) {
			r.Use(authRequired)
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateMe)
			r.Delete("/me", userHandler.DeleteMe)
			r.Put("/me/preferences", userHandler.UpdatePreferences)
		})

		r.Route("/auth/2fa", func(r chi.Router) {
			r.Use(authRequired)
			r.Post("/setup", authHandler.Setup2FA)
			r.Post("/verify", authHandler.Verify2FA)
			r.Post("/disable", authHandler.Disable2FA)
		})

		r.Route("/matches", func(r chi.Router) {
			r.Use(authRequired)
			r.Get("/candidates", matchHandler.GetCandidates)
			r.Post("/swipe", matchHandler.Swipe)
			r.Get("/", matchHandler.GetMatches)
			r.Get("/{id}", matchHandler.GetMatch)
		})
	})

	// 9. Start server with graceful shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info().Msgf("server listening on :%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("server forced to shutdown")
	}
	logger.Info().Msg("server stopped")
}
