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

// tokenClaims 是 access/refresh 共用的 JWT 载荷结构。
// Use 字段用于显式区分 access 和 refresh，防止跨用途误用。
type tokenClaims struct {
	UserID       uint64 `json:"uid"`
	SID          string `json:"sid"`
	TokenVersion uint32 `json:"tv"`
	Use          string `json:"use"`
	jwt.RegisteredClaims
}

// newSID 生成登录会话 SID。
func (s *Service) newSID() (string, error) {
	return randomToken(10)
}

// issueTokenPair 签发 access+refresh 令牌对，并返回 refresh token 哈希。
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

	// access token 用于网关和业务接口鉴权，生命周期更短。
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
	// refresh token 仅用于换发新 token，对应 refresh_session 持久化状态。
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

// parseToken 统一执行 JWT 解析与基础合法性校验。
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

// parseAccessToken 校验并解析 access token。
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

// parseRefreshToken 校验并解析 refresh token。
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
