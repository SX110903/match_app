//go:build ignore

package main

import (
	"fmt"

	"golang.org/x/crypto/argon2"
	"crypto/rand"
	"encoding/base64"
)

func main() {
	password := "Admin@MatchHub2026!SuperSecure#Xk9mP3vLqR7nZwYtFbDcHjGs5eUo1iA"

	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=65536,t=1,p=4$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	fmt.Println(encoded)
}
