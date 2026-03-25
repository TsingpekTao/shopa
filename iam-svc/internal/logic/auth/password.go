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
	// Argon2id 参数：在线鉴权场景下优先保证抗暴力破解能力，同时兼顾服务延迟。
	// 这些参数需与 verifyPassword 使用同一组配置，否则会导致历史密码无法通过校验。
	argon2Time    uint32 = 3
	argon2Memory  uint32 = 64 * 1024
	argon2Threads uint8  = 2
	argon2KeyLen  uint32 = 32
	saltLen              = 16
)

// makePasswordHash 生成密码哈希。
// 返回值约定：
// 1) `hash`/`salt` 采用 base64.RawStdEncoding，便于数据库存储且不含填充符。
// 2) `algo`/`version` 用于后续算法平滑升级（例如切换参数或算法时做兼容验证）。
func makePasswordHash(password string) (hash, salt, algo string, version uint, err error) {
	// 每次都生成独立随机盐，防止相同密码产出相同哈希。
	rawSalt := make([]byte, saltLen)
	if _, err = rand.Read(rawSalt); err != nil {
		return "", "", "", 0, err
	}
	key := argon2.IDKey([]byte(password), rawSalt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(key), base64.RawStdEncoding.EncodeToString(rawSalt), "argon2id", 1, nil
}

// verifyPassword 校验明文密码是否与已存哈希匹配。
// 安全要点：
// 1) 先做空值/算法白名单校验，避免异常数据绕过。
// 2) 最终比较使用 ConstantTimeCompare，降低时序侧信道风险。
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

// validateNewPassword 校验新密码复杂度。
// 规则：8~20 位，且必须同时包含字母、数字、特殊字符；不允许空白符。
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
