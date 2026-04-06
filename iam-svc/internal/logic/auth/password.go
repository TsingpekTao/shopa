package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    uint32 = 3
	argon2Memory  uint32 = 64 * 1024
	argon2Threads uint8  = 2
	argon2KeyLen  uint32 = 32
	saltLen              = 16
)

func makePasswordHash(password string) (hash, salt, algo string, version uint, err error) {
	rawSalt := make([]byte, saltLen)
	if _, err = rand.Read(rawSalt); err != nil {
		return "", "", "", 0, err
	}
	key := argon2.IDKey([]byte(password), rawSalt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(key), base64.RawStdEncoding.EncodeToString(rawSalt), "argon2id", 1, nil
}

func verifyPassword(password, hash, salt, algo string) bool {
	if strings.TrimSpace(hash) == "" || strings.TrimSpace(salt) == "" {
		return false
	}
	if algo != "" && !strings.EqualFold(algo, "argon2id") {
		return false
	}

	rawSalt, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}
	rawHash, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	key := argon2.IDKey([]byte(password), rawSalt, argon2Time, argon2Memory, argon2Threads, uint32(len(rawHash)))
	return subtle.ConstantTimeCompare(rawHash, key) == 1
}

func validateNewPassword(p string) error {
	password := p
	length := len([]rune(password))
	if length < 8 || length > 20 {
		return fmt.Errorf("\u5bc6\u7801\u957f\u5ea6\u5fc5\u987b\u4e3a 8-20 \u4f4d")
	}

	var (
		hasLetter  bool
		hasDigit   bool
		hasSpecial bool
	)
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r), unicode.IsSymbol(r):
			hasSpecial = true
		case unicode.IsSpace(r):
			return fmt.Errorf("\u5bc6\u7801\u4e0d\u80fd\u5305\u542b\u7a7a\u683c")
		}
	}

	if !hasLetter || !hasDigit || !hasSpecial {
		return fmt.Errorf("\u5bc6\u7801\u5fc5\u987b\u540c\u65f6\u5305\u542b\u5b57\u6bcd\u3001\u6570\u5b57\u548c\u7279\u6b8a\u5b57\u7b26")
	}
	return nil
}
