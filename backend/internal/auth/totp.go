package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type ITOTPService interface {
	GenerateSecret(email string) (*otp.Key, error)
	Validate(code, secret string) bool
	GenerateBackupCodes(n int) ([]string, error)
}

type totpService struct{}

func NewTOTPService() ITOTPService {
	return &totpService{}
}

func (s *totpService) GenerateSecret(email string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MatchHub",
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("generating TOTP secret: %w", err)
	}
	return key, nil
}

func (s *totpService) Validate(code, secret string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes produces n one-time-use codes as hex strings.
func (s *totpService) GenerateBackupCodes(n int) ([]string, error) {
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 5)
		if _, err := rand.Read(b); err != nil {
			return nil, fmt.Errorf("generating backup code: %w", err)
		}
		codes[i] = hex.EncodeToString(b)
	}
	return codes, nil
}
