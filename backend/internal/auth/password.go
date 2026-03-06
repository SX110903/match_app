package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters (2024 recommendations)
const (
	argon2Memory      = 64 * 1024 // 64MB
	argon2Iterations  = 3
	argon2Parallelism = 4
	argon2SaltLength  = 16
	argon2KeyLength   = 32
)

// HashPassword creates an Argon2id hash of the password.
// Returns the encoded string in the format: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
func HashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLength)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Memory,
		argon2Iterations,
		argon2Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks a password against an Argon2id hash using constant-time comparison.
func VerifyPassword(password, encodedHash string) (bool, error) {
	var version int
	var memory, iterations uint32
	var parallelism uint8
	var saltB64, hashB64 string

	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		&version, &memory, &iterations, &parallelism, &saltB64,
	)
	if err != nil {
		// Try splitting manually for the last two parts
		parts := splitArgon2Hash(encodedHash)
		if len(parts) != 6 {
			return false, fmt.Errorf("invalid hash format")
		}
		saltB64 = parts[4]
		hashB64 = parts[5]
	} else {
		parts := splitArgon2Hash(encodedHash)
		if len(parts) != 6 {
			return false, fmt.Errorf("invalid hash format")
		}
		hashB64 = parts[5]
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, fmt.Errorf("decoding salt: %w", err)
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, fmt.Errorf("decoding hash: %w", err)
	}

	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(storedHash)))

	return subtle.ConstantTimeCompare(storedHash, computedHash) == 1, nil
}

func splitArgon2Hash(s string) []string {
	parts := make([]string, 0, 6)
	current := ""
	count := 0
	for _, c := range s {
		if c == '$' {
			if count > 0 {
				parts = append(parts, current)
				current = ""
			}
			count++
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
