package email

import "context"

type IEmailService interface {
	SendVerificationEmail(ctx context.Context, to, name, verifyURL string) error
	SendPasswordResetEmail(ctx context.Context, to, resetURL string) error
	SendPasswordChangedEmail(ctx context.Context, to string) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
}
