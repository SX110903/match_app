package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/SX110903/match_app/backend/internal/config"
)

type Claims struct {
	jwt.RegisteredClaims
	Roles       []string `json:"roles"`
	TwoFAVerified bool   `json:"2fa_verified"`
	TokenType   string   `json:"token_type,omitempty"` // "temp_2fa" for 2FA temp tokens
}

type IJWTService interface {
	GenerateAccessToken(userID string, roles []string, twoFAVerified bool) (string, string, error) // token, jti, error
	GenerateRefreshToken(userID string) (string, error)
	GenerateTempToken(userID string) (string, error) // For 2FA second step
	ValidateAccessToken(tokenStr string) (*Claims, error)
	ValidateTempToken(tokenStr string) (*Claims, error)
}

type jwtService struct {
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(cfg config.JWTConfig) (IJWTService, error) {
	privateKey, err := loadPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("loading private key: %w", err)
	}
	publicKey, err := loadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("loading public key: %w", err)
	}
	return &jwtService{
		privateKey:    privateKey,
		publicKey:     publicKey,
		accessExpiry:  cfg.AccessTokenExpiry,
		refreshExpiry: cfg.RefreshTokenExpiry,
	}, nil
}

func (s *jwtService) GenerateAccessToken(userID string, roles []string, twoFAVerified bool) (string, string, error) {
	jti := uuid.New().String()
	now := time.Now()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
		},
		Roles:         roles,
		TwoFAVerified: twoFAVerified,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("signing access token: %w", err)
	}
	return signed, jti, nil
}

func (s *jwtService) GenerateRefreshToken(userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating refresh token: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *jwtService) GenerateTempToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
		},
		TokenType: "temp_2fa",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *jwtService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	return s.validateToken(tokenStr, "")
}

func (s *jwtService) ValidateTempToken(tokenStr string) (*Claims, error) {
	return s.validateToken(tokenStr, "temp_2fa")
}

func (s *jwtService) validateToken(tokenStr, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if expectedType != "" && claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

var ErrInvalidToken = fmt.Errorf("invalid token")

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return rsaKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaKey, nil
}

// Context key for user claims
type contextKey string
const ClaimsKey contextKey = "claims"

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*Claims)
	return claims, ok
}

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}
