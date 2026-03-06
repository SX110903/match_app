package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/internal/validator"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type AuthHandler struct {
	authSvc service.IAuthService
}

func NewAuthHandler(authSvc service.IAuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	if err := h.authSvc.Register(r.Context(), req); err != nil {
		switch err {
		case domain.ErrConflict:
			response.Conflict(w, "email already in use")
		default:
			logger.Error().Err(err).Msg("register failed")
			response.InternalError(w)
		}
		return
	}

	response.Created(w, map[string]string{"message": "registration successful, please verify your email"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	result, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			response.Unauthorized(w, "invalid credentials")
		case domain.ErrEmailNotVerified:
			response.Unauthorized(w, "email not verified")
		case domain.ErrAccountLocked:
			response.Error(w, http.StatusTooManyRequests, "account temporarily locked")
		default:
			logger.Error().Err(err).Msg("login failed")
			response.InternalError(w)
		}
		return
	}

	// If 2FA required, don't set cookie yet
	if result.Requires2FA {
		response.OK(w, result)
		return
	}

	// Only access_token goes in response body (stays in memory on client)
	response.OK(w, map[string]string{"access_token": result.AccessToken})
}

func (h *AuthHandler) LoginWith2FA(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TempToken string `json:"temp_token" validate:"required"`
		Code      string `json:"code"       validate:"required,len=6"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	result, err := h.authSvc.LoginWith2FA(r.Context(), req.TempToken, req.Code)
	if err != nil {
		switch err {
		case domain.ErrTokenInvalid:
			response.Unauthorized(w, "invalid or expired token")
		case domain.ErrTwoFAInvalid:
			response.Unauthorized(w, "invalid 2FA code")
		default:
			response.InternalError(w)
		}
		return
	}

	response.OK(w, map[string]string{"access_token": result.AccessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	refreshToken, _ := r.Cookie("refresh_token")
	rt := ""
	if refreshToken != nil {
		rt = refreshToken.Value
	}

	expiry := int64(0)
	if claims.ExpiresAt != nil {
		expiry = claims.ExpiresAt.Unix()
	}

	if err := h.authSvc.Logout(r.Context(), claims.ID, rt, expiry); err != nil {
		logger.Error().Err(err).Msg("logout failed")
		response.InternalError(w)
		return
	}

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/auth",
	})

	response.OK(w, map[string]string{"message": "logged out successfully"})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Unauthorized(w, "refresh token not found")
		return
	}

	result, err := h.authSvc.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		response.Unauthorized(w, "invalid or expired refresh token")
		return
	}

	response.OK(w, map[string]string{"access_token": result.AccessToken})
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.authSvc.VerifyEmail(r.Context(), req.Token); err != nil {
		response.BadRequest(w, "invalid or expired verification token")
		return
	}

	response.OK(w, map[string]string{"message": "email verified successfully"})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Always same response to prevent user enumeration
	_ = h.authSvc.ForgotPassword(r.Context(), req.Email)
	response.OK(w, map[string]string{"message": "if the email exists, you will receive a reset link"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"        validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8,max=128"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	if err := h.authSvc.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		response.BadRequest(w, "invalid or expired reset token")
		return
	}

	response.OK(w, map[string]string{"message": "password reset successfully"})
}

func (h *AuthHandler) Setup2FA(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	result, err := h.authSvc.Setup2FA(r.Context(), claims.Subject)
	if err != nil {
		switch err {
		case domain.ErrTwoFAAlreadyEnabled:
			response.Conflict(w, "2FA is already enabled")
		default:
			logger.Error().Err(err).Msg("setup 2fa failed")
			response.InternalError(w)
		}
		return
	}

	response.OK(w, result)
}

func (h *AuthHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	var req struct {
		Code string `json:"code" validate:"required,len=6"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.authSvc.Verify2FA(r.Context(), claims.Subject, req.Code); err != nil {
		response.Unauthorized(w, "invalid 2FA code")
		return
	}

	response.OK(w, map[string]string{"message": "2FA enabled successfully"})
}

func (h *AuthHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	var req struct {
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.authSvc.Disable2FA(r.Context(), claims.Subject, req.Password); err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			response.Unauthorized(w, "invalid password")
		case domain.ErrTwoFANotEnabled:
			response.BadRequest(w, "2FA is not enabled")
		default:
			response.InternalError(w)
		}
		return
	}

	response.OK(w, map[string]string{"message": "2FA disabled successfully"})
}
