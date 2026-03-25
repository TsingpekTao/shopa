package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"github.com/golang-jwt/jwt/v5"
)

// tokenClaims 是 access/refresh 共享的 JWT 负载结构。
// `Use` 字段是关键安全标记：同一个签名密钥下，必须显式区分 access 与 refresh 用途，
// 防止“拿 refresh token 调 access 接口”一类跨用途误用。
type tokenClaims struct {
	UserID       uint64 `json:"uid"`
	SID          string `json:"sid"`
	TokenVersion uint32 `json:"tv"`
	Use          string `json:"use"`
	jwt.RegisteredClaims
}

// newSID 生成登录会话 SID。
// access 与 refresh 会共享同一个 SID，便于追踪会话链路与执行会话级注销。
func (s *Service) newSID() (string, error) {
	return randomToken(10)
}

// issueTokenPair 签发 access+refresh 成对令牌。
// 返回 refreshHash 的目的：数据库只落 hash，不存 refresh 明文，降低泄露风险。
func (s *Service) issueTokenPair(userID uint64, sid string, tokenVersion uint32) (
	pair *v1.TokenPair,
	refreshHash string,
	accessExp time.Time,
	refreshExp time.Time,
	err error,
) {
	var (
		nowTime   = now()
		issuer    = "shopa-iam"
		sub       = strconv.FormatUint(userID, 10)
		accessTTL = s.accessTTL()
		refrTTL   = s.refreshTTL()
	)

	accessExp = nowTime.Add(accessTTL)
	refreshExp = nowTime.Add(refrTTL)

	// access token：用于网关和业务接口鉴权，生命周期更短。
	accessClaims := tokenClaims{
		UserID:       userID,
		SID:          sid,
		TokenVersion: tokenVersion,
		Use:          consts.TokenUseAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}
	// refresh token：仅用于换发新 token，对应 refresh_session 持久化状态。
	refreshClaims := tokenClaims{
		UserID:       userID,
		SID:          sid,
		TokenVersion: tokenVersion,
		Use:          consts.TokenUseRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			ExpiresAt: jwt.NewNumericDate(refreshExp),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwt.Secret))
	if err != nil {
		return nil, "", time.Time{}, time.Time{}, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.jwt.Secret))
	if err != nil {
		return nil, "", time.Time{}, time.Time{}, err
	}

	refreshHash = sha256Hex(refreshToken)
	pair = &v1.TokenPair{
		TokenType:        consts.TokenTypeBearer,
		AccessToken:      accessToken,
		AccessExpiresIn:  uint32(accessTTL / time.Second),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: uint32(refrTTL / time.Second),
		Sid:              sid,
	}
	return pair, refreshHash, accessExp, refreshExp, nil
}

// parseToken 统一执行 JWT 解析与基础校验。
// 校验项：
// 1) token 非空。
// 2) 签名算法必须是 HMAC，避免算法混淆攻击。
// 3) 签名与时效（exp/nbf 等）由 jwt 库校验。
func (s *Service) parseToken(token string) (*tokenClaims, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errs.New(errs.CodeAccessTokenInvalid)
	}
	claims := &tokenClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.jwt.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.New(errs.CodeAccessTokenExpired)
		}
		return nil, errs.Wrap(errs.CodeAccessTokenInvalid, err)
	}
	if parsed == nil || !parsed.Valid {
		return nil, errs.New(errs.CodeAccessTokenInvalid)
	}
	return claims, nil
}

// parseAccessToken 解析并校验 access token（use=access）。
func (s *Service) parseAccessToken(token string) (*tokenClaims, error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return nil, err
	}
	if claims.Use != consts.TokenUseAccess {
		return nil, errs.New(errs.CodeAccessTokenInvalid)
	}
	return claims, nil
}

// parseRefreshToken 解析并校验 refresh token（use=refresh）。
// 与 access 分离错误码，便于上层区分“刷新失败”与“访问失败”。
func (s *Service) parseRefreshToken(token string) (*tokenClaims, error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}
	if claims.Use != consts.TokenUseRefresh {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}
	return claims, nil
}
