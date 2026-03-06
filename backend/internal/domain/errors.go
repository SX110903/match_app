package domain

import "errors"

// Domain errors - typed, never exposed raw to the client
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrConflict          = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified  = errors.New("email not verified")
	ErrAccountLocked     = errors.New("account temporarily locked")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrTwoFARequired     = errors.New("two-factor authentication required")
	ErrTwoFAInvalid      = errors.New("invalid two-factor code")
	ErrTwoFANotEnabled   = errors.New("two-factor authentication not enabled")
	ErrTwoFAAlreadyEnabled = errors.New("two-factor authentication already enabled")
	ErrWeakPassword      = errors.New("password does not meet requirements")
	ErrSelfAction        = errors.New("cannot perform this action on yourself")
)
