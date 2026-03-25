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

// makePasswordHash 生成 Argon2id 哈希与随机盐。
func makePasswordHash(password string) (hash, salt, algo string, version uint, err error) {
	rawSalt := make([]byte, saltLen)
	if _, err = rand.Read(rawSalt); err != nil {
		return "", "", "", 0, err
	}
	key := argon2.IDKey([]byte(password), rawSalt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(key), base64.RawStdEncoding.EncodeToString(rawSalt), "argon2id", 1, nil
}

// verifyPassword 校验明文密码与存储哈希是否匹配。
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

// validateNewPassword 校验密码复杂度：8-20 位，且必须包含字母、数字、特殊字符。
func validateNewPassword(p string) error {
	password := strings.TrimSpace(p)
	length := len([]rune(password))
	if length < 8 || length > 20 {
		return fmt.Errorf("密码长度必须在 8 到 20 位之间")
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
			return fmt.Errorf("密码不能包含空格")
		}
	}

	if !hasLetter || !hasDigit || !hasSpecial {
		return fmt.Errorf("密码必须同时包含字母、数字和特殊字符")
	}
	return nil
}
