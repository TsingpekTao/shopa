package auth

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/dao"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/sms"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
)

// SendSmsCode 发送短信验证码并返回重发窗口。
// 关键路径：参数校验 -> 风控校验 -> 分布式发送锁 -> 生成验证码 -> 缓存存储哈希 -> 审计日志。
func (s *Service) SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (*v1.SendSmsCodeRes, error) {
	var (
		scene = req.GetScene()
		phone = normalizePhone(req.GetPhone())
		meta  = extractRiskMeta(ctx, req.GetRisk())
	)

	// 场景和手机号是最小输入集，缺失直接返回参数错误。
	if scene == v1.SmsScene_SMS_SCENE_UNSPECIFIED || phone == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	// 轻量自动化识别命中时要求 captcha，降低机器批量探测风险。
	if looksLikeAutomation(meta.UserAgent) && meta.CaptchaToken == "" {
		return nil, errs.New(errs.CodeCaptchaFailed)
	}

	// 注册场景附加 IP 频控，抑制同源批量注册。
	if scene == v1.SmsScene_SMS_SCENE_REGISTER && meta.ClientIP != "" {
		count, err := s.cache.IncRegisterIpLimit(ctx, meta.ClientIP, time.Now())
		if err != nil {
			return nil, errs.Wrap(errs.CodeInternalError, err)
		}
		if count > 3 {
			return nil, errs.New(errs.CodeRegisterTooFrequent)
		}
	}

	// 分布式发送锁：同手机号+场景短时间内只允许一个发送请求成功。
	ok, err := s.cache.AcquireSmsSendLock(ctx, sceneKey(scene), phone)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if !ok {
		return nil, errs.New(errs.CodeSmsTooFrequent)
	}

	// 生成纯数字验证码（当前长度由常量统一控制）。
	code, err := randomDigits(smsCodeLength)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	// 缓存中仅保存验证码哈希，不保存明文。
	if err = s.cache.SaveSmsCode(ctx, sceneKey(scene), phone, s.hashSmsCode(scene, phone, code)); err != nil {
		s.insertSmsLog(ctx, scene, phone, consts.SmsProviderMock, "", false, errs.CodeInternalError.Code(), meta)
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	// 通过 provider 发送验证码。若发送失败，回滚缓存验证码，避免产生“用户收不到但服务端可校验”的幽灵验证码。
	sendRes, err := s.sms.SendCode(ctx, &sms.SendCodeRequest{
		Scene:     sceneKey(scene),
		Phone:     phone,
		Code:      code,
		RequestID: meta.RequestID,
		ClientIP:  meta.ClientIP,
	})
	if err != nil {
		_ = s.cache.DeleteSmsCode(ctx, sceneKey(scene), phone)
		provider, bizID := consts.SmsProviderMock, ""
		if sendRes != nil {
			if strings.TrimSpace(sendRes.Provider) != "" {
				provider = sendRes.Provider
			}
			bizID = strings.TrimSpace(sendRes.BizID)
		}
		s.insertSmsLog(ctx, scene, phone, provider, bizID, false, errs.CodeInternalError.Code(), meta)
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	// 发送成功后记录 provider/biz_id，便于后续对账与排障。
	provider, bizID := consts.SmsProviderMock, ""
	if sendRes != nil {
		if strings.TrimSpace(sendRes.Provider) != "" {
			provider = sendRes.Provider
		}
		bizID = strings.TrimSpace(sendRes.BizID)
	}
	// 审计日志写入失败不阻断主流程（insertSmsLog 内部已忽略写库错误）。
	s.insertSmsLog(ctx, scene, phone, provider, bizID, true, errs.CodeOK.Code(), meta)

	return &v1.SendSmsCodeRes{ResendAfterSeconds: 60}, nil
}

// RegisterByPassword 使用“手机号 + 短信码 + 密码”完成注册。
// 关键事务：用户主表、默认角色、默认会员、refresh_session、outbox 事件同事务提交。
func (s *Service) RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (*v1.RegisterByPasswordRes, error) {
	var (
		phone    = normalizePhone(req.GetPhone())
		smsCode  = strings.TrimSpace(req.GetSmsCode())
		password = req.GetPassword()
		meta     = extractRiskMeta(ctx, req.GetRisk())
	)

	if phone == "" || smsCode == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}
	if err := validateNewPassword(password); err != nil {
		return nil, errs.New(errs.CodeInvalidParam, err.Error())
	}

	if looksLikeAutomation(meta.UserAgent) && meta.CaptchaToken == "" {
		return nil, errs.New(errs.CodeCaptchaFailed)
	}

	// 注册链路单独做 IP 频控，阻断同 IP 快速重复注册。
	if meta.ClientIP != "" {
		count, err := s.cache.IncRegisterIpLimit(ctx, meta.ClientIP, time.Now())
		if err != nil {
			return nil, errs.Wrap(errs.CodeInternalError, err)
		}
		if count > 3 {
			return nil, errs.New(errs.CodeRegisterTooFrequent)
		}
	}

	if err := s.verifySmsCode(ctx, v1.SmsScene_SMS_SCENE_REGISTER, phone, smsCode); err != nil {
		return nil, err
	}

	// 注册前先查重，减少无效事务开销。
	exist, err := s.findAuthByPhone(ctx, phone)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if exist != nil {
		return nil, errs.New(errs.CodePhoneRegistered)
	}

	hash, salt, algo, ver, err := makePasswordHash(password)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	initName, err := s.makeInitDisplayName()
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	var (
		userID = s.nextUserID()
		// tokenVersion 作为令牌批次号，后续改密/封禁可通过递增实现全量失效。
		tokenVersion = uint(1)
	)
	// sid 代表一次登录会话，access/refresh 共用同一个 sid。
	sid, sidErr := s.newSID()
	if sidErr != nil {
		return nil, errs.Wrap(errs.CodeInternalError, sidErr)
	}
	pair, refreshHash, _, refreshExp, err := s.issueTokenPair(userID, sid, uint32(tokenVersion))
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	nowTime := gtime.NewFromTime(time.Now().UTC())
	err = dao.IamUserAuth.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 事务步骤 1：写入认证主记录。
		if _, err = tx.Model(dao.IamUserAuth.Table()).Data(do.IamUserAuth{
			UserId:               userID,
			Phone:                phone,
			PasswordHash:         hash,
			PasswordSalt:         salt,
			PasswordAlgo:         algo,
			PasswordVer:          ver,
			AccountStatus:        consts.AccountStatusActive,
			FailedLoginCount:     0,
			TokenVersion:         tokenVersion,
			LastLoginAt:          nowTime,
			LastLoginIp:          meta.ClientIP,
			LastLoginGeo:         "",
			LastLoginUa:          meta.UserAgent,
			LastLoginFingerprint: meta.Fingerprint,
		}).Insert(); err != nil {
			return err
		}

		// 事务步骤 2：分配默认角色。
		if _, err = tx.Model(dao.IamUserRole.Table()).Data(do.IamUserRole{
			UserId:    userID,
			RoleCode:  consts.RoleCodeCustomer,
			ScopeType: consts.ScopeTypeGlobal,
			ScopeId:   0,
			Status:    consts.StatusActive,
		}).Insert(); err != nil {
			return err
		}

		// 事务步骤 3：初始化会员信息。
		if _, err = tx.Model(dao.IamMembership.Table()).Data(do.IamMembership{
			UserId:    userID,
			LevelCode: consts.MembershipLevelBasic,
			Points:    0,
		}).Insert(); err != nil {
			return err
		}

		// 事务步骤 4：写入 refresh session。
		if _, err = tx.Model(dao.IamRefreshSession.Table()).Data(do.IamRefreshSession{
			Sid:              sid,
			UserId:           userID,
			RefreshTokenHash: refreshHash,
			UaHash:           sha256Hex(meta.UserAgent),
			Ip:               meta.ClientIP,
			Geo:              "",
			Fingerprint:      meta.Fingerprint,
			ExpiresAt:        gtime.NewFromTime(refreshExp.UTC()),
		}).Insert(); err != nil {
			return err
		}

		// 事务步骤 5：写 outbox，保证“注册成功”与“事件可投递”原子一致。
		return s.insertOutboxUserRegistered(ctx, tx, userID, initName)
	})
	if err != nil {
		// 并发注册触发唯一键冲突时，转换为业务错误码。
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, errs.New(errs.CodePhoneRegistered)
		}
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	auth, err := s.findAuthByUserID(ctx, userID)
	if err != nil || auth == nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	session, err := s.buildSessionSummary(ctx, auth)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.RegisterByPasswordRes{
		UserId:          userID,
		InitDisplayName: initName,
		Auth:            s.buildTokenAuthResult(pair),
		Session:         session,
	}, nil
}

// LoginByPassword 使用账号密码登录。
// 关键分支：锁定校验 -> 账号状态校验 -> 密码校验 -> 可选 MFA -> 登录收尾。
func (s *Service) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginByPasswordRes, error) {
	var (
		identifier = normalizeIdentifier(req.GetIdentifier())
		password   = req.GetPassword()
		meta       = extractRiskMeta(ctx, req.GetRisk())
	)
	if identifier == "" || strings.TrimSpace(password) == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	// 先查 Redis 级锁定（短路返回，避免打数据库）。
	if locked, err := s.cache.IsLoginLocked(ctx, identifier); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	} else if locked {
		s.insertLoginLog(ctx, 0, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, false, "redis_lock", meta)
		return nil, errs.New(errs.CodeAccountLocked)
	}

	auth, err := s.findAuthByIdentifier(ctx, identifier)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if auth == nil {
		return nil, s.onLoginFailed(ctx, nil, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, meta, "user_not_found")
	}

	if auth.AccountStatus == consts.AccountStatusDisabled {
		s.insertLoginLog(ctx, auth.UserId, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, false, "disabled", meta)
		return nil, errs.New(errs.CodeAccountDisabled)
	}
	if auth.AccountStatus == consts.AccountStatusLocked && auth.LockedUntil != nil && auth.LockedUntil.Timestamp() > time.Now().UTC().Unix() {
		s.insertLoginLog(ctx, auth.UserId, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, false, "locked_until", meta)
		return nil, errs.New(errs.CodeAccountLocked)
	}

	if !verifyPassword(password, auth.PasswordHash, auth.PasswordSalt, auth.PasswordAlgo) {
		return nil, s.onLoginFailed(ctx, auth, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, meta, "invalid_password")
	}

	// 策略命中：IP 变化触发 MFA 二次验证，不直接签发 token。
	if s.security.MfaOnIPChange && shouldRequireMFA(auth, meta.ClientIP) {
		_, challenge, err := s.buildMFAChallenge(ctx, auth, identifier)
		if err != nil {
			return nil, err
		}
		s.insertLoginLog(ctx, auth.UserId, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, false, "mfa_required", meta)
		return &v1.LoginByPasswordRes{
			Channel: v1.LoginChannel_LOGIN_CHANNEL_PASSWORD,
			Auth:    challenge,
			Session: &v1.SessionSummary{UserId: auth.UserId},
		}, nil
	}

	pair, session, err := s.finishLogin(ctx, auth, identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, meta)
	if err != nil {
		return nil, err
	}
	return &v1.LoginByPasswordRes{
		Channel: v1.LoginChannel_LOGIN_CHANNEL_PASSWORD,
		Auth:    s.buildTokenAuthResult(pair),
		Session: session,
	}, nil
}

// LoginBySms 使用手机号 + 短信验证码登录。
func (s *Service) LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (*v1.LoginBySmsRes, error) {
	var (
		phone   = normalizePhone(req.GetPhone())
		smsCode = strings.TrimSpace(req.GetSmsCode())
		meta    = extractRiskMeta(ctx, req.GetRisk())
	)
	if phone == "" || smsCode == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	if locked, err := s.cache.IsLoginLocked(ctx, phone); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	} else if locked {
		s.insertLoginLog(ctx, 0, phone, v1.LoginChannel_LOGIN_CHANNEL_SMS, false, "redis_lock", meta)
		return nil, errs.New(errs.CodeAccountLocked)
	}

	if err := s.verifySmsCode(ctx, v1.SmsScene_SMS_SCENE_LOGIN, phone, smsCode); err != nil {
		s.insertLoginLog(ctx, 0, phone, v1.LoginChannel_LOGIN_CHANNEL_SMS, false, "sms_verify_failed", meta)
		return nil, err
	}

	auth, err := s.findAuthByPhone(ctx, phone)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if auth == nil {
		return nil, s.onLoginFailed(ctx, nil, phone, v1.LoginChannel_LOGIN_CHANNEL_SMS, meta, "user_not_found")
	}

	if auth.AccountStatus == consts.AccountStatusDisabled {
		s.insertLoginLog(ctx, auth.UserId, phone, v1.LoginChannel_LOGIN_CHANNEL_SMS, false, "disabled", meta)
		return nil, errs.New(errs.CodeAccountDisabled)
	}
	if auth.AccountStatus == consts.AccountStatusLocked && auth.LockedUntil != nil && auth.LockedUntil.Timestamp() > time.Now().UTC().Unix() {
		s.insertLoginLog(ctx, auth.UserId, phone, v1.LoginChannel_LOGIN_CHANNEL_SMS, false, "locked_until", meta)
		return nil, errs.New(errs.CodeAccountLocked)
	}

	pair, session, err := s.finishLogin(ctx, auth, phone, v1.LoginChannel_LOGIN_CHANNEL_SMS, meta)
	if err != nil {
		return nil, err
	}
	return &v1.LoginBySmsRes{
		Channel: v1.LoginChannel_LOGIN_CHANNEL_SMS,
		Auth:    s.buildTokenAuthResult(pair),
		Session: session,
	}, nil
}

// VerifyMfaChallenge 校验 challenge + 短信码，成功后补发 token。
func (s *Service) VerifyMfaChallenge(ctx context.Context, req *v1.VerifyMfaChallengeReq) (*v1.VerifyMfaChallengeRes, error) {
	var (
		challengeID = strings.TrimSpace(req.GetChallengeId())
		code        = strings.TrimSpace(req.GetSmsCode())
		meta        = extractRiskMeta(ctx, req.GetRisk())
	)
	if challengeID == "" || code == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	raw, exists, err := s.cache.GetMFAChallenge(ctx, challengeID)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if !exists {
		return nil, errs.New(errs.CodeMFAFailed)
	}

	var payload mfaChallengePayload
	if err = json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, errs.New(errs.CodeMFAFailed)
	}
	if payload.UserID == 0 || payload.Phone == "" {
		return nil, errs.New(errs.CodeMFAFailed)
	}

	if err = s.verifySmsCode(ctx, v1.SmsScene_SMS_SCENE_MFA, payload.Phone, code); err != nil {
		return nil, errs.New(errs.CodeMFAFailed)
	}

	auth, err := s.findAuthByUserID(ctx, payload.UserID)
	if err != nil || auth == nil {
		return nil, errs.New(errs.CodeInvalidCredential)
	}

	pair, session, err := s.finishLogin(ctx, auth, payload.Identifier, v1.LoginChannel_LOGIN_CHANNEL_PASSWORD, meta)
	if err != nil {
		return nil, err
	}
	// challenge 一次性消费，防止重复使用。
	_ = s.cache.DeleteMFAChallenge(ctx, challengeID)

	return &v1.VerifyMfaChallengeRes{
		TokenPair: pair,
		Session:   session,
	}, nil
}

// verifySmsCode 校验短信码，并处理错误次数上限。
// 达到阈值后主动删除验证码，阻断无限试错。
func (s *Service) verifySmsCode(ctx context.Context, scene v1.SmsScene, phone, code string) error {
	sceneName := sceneKey(scene)
	hash, exists, err := s.cache.GetSmsCodeHash(ctx, sceneName, phone)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, err)
	}
	if !exists {
		return errs.New(errs.CodeSmsCodeExpired)
	}
	if hash != s.hashSmsCode(scene, phone, code) {
		attempt, _ := s.cache.IncrSmsCodeAttempt(ctx, sceneName, phone)
		if attempt >= smsCodeMaxAttempts {
			_ = s.cache.DeleteSmsCode(ctx, sceneName, phone)
		}
		return errs.New(errs.CodeSmsCodeInvalid)
	}
	_ = s.cache.DeleteSmsCode(ctx, sceneName, phone)
	return nil
}

// onLoginFailed 统一处理登录失败分支。
// 执行动作：
// 1) 增加失败计数并施加延迟。
// 2) 按阈值锁定账号（数据库状态 + Redis 快速锁）。
// 3) 写审计日志。
func (s *Service) onLoginFailed(ctx context.Context, auth *entity.IamUserAuth, identifier string, channel v1.LoginChannel, meta riskMeta, reason string) error {
	// count 表示连续失败次数，用于锁定判断与退避延迟计算。
	count, err := s.cache.IncrLoginFail(ctx, identifier)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, err)
	}
	sleepOnFailure(count)

	if auth != nil {
		cols := dao.IamUserAuth.Columns()
		data := do.IamUserAuth{
			FailedLoginCount: gdb.Raw(cols.FailedLoginCount + " + 1"),
		}
		if count >= s.security.LoginFailMax {
			data.AccountStatus = consts.AccountStatusLocked
			data.LockedUntil = gtime.NewFromTime(time.Now().UTC().Add(consts.DefaultLockDuration))
		}
		_, _ = dao.IamUserAuth.Ctx(ctx).Where(cols.UserId, auth.UserId).Data(data).Update()
	}

	if count >= s.security.LoginFailMax {
		_ = s.cache.LockLogin(ctx, identifier)
	}

	var uid uint64
	if auth != nil {
		uid = auth.UserId
	}
	s.insertLoginLog(ctx, uid, identifier, channel, false, reason, meta)

	if count >= s.security.LoginFailMax {
		return errs.New(errs.CodeAccountLocked)
	}
	return errs.New(errs.CodeInvalidCredential)
}

// buildMFAChallenge 生成 MFA challenge 并写入 Redis。
// VerifyMfaChallenge 会基于该 challenge 做二次校验并发放 token。
func (s *Service) buildMFAChallenge(ctx context.Context, auth *entity.IamUserAuth, identifier string) (string, *v1.AuthResult, error) {
	id, err := randomToken(10)
	if err != nil {
		return "", nil, errs.Wrap(errs.CodeInternalError, err)
	}

	payload, _ := json.Marshal(mfaChallengePayload{
		UserID:     auth.UserId,
		Phone:      auth.Phone,
		Identifier: identifier,
	})
	if err = s.cache.SaveMFAChallenge(ctx, id, string(payload)); err != nil {
		return "", nil, errs.Wrap(errs.CodeInternalError, err)
	}

	exp := time.Now().UTC().Add(consts.DefaultMfaTTL)
	return id, s.buildMFAAuthResult(id, exp), nil
}

// shouldRequireMFA 判断是否触发异地登录 MFA（当前按 IP 变化判定）。
func shouldRequireMFA(auth *entity.IamUserAuth, currentIP string) bool {
	if strings.TrimSpace(currentIP) == "" {
		return false
	}
	if strings.TrimSpace(auth.LastLoginIp) == "" {
		return false
	}
	return strings.TrimSpace(auth.LastLoginIp) != strings.TrimSpace(currentIP)
}

// finishLogin 处理登录成功收尾。
// 关键顺序：
// 1) 重置失败计数。
// 2) 签发 token 并落 refresh_session。
// 3) 更新最近登录信息。
// 4) 查询并返回会话摘要。
func (s *Service) finishLogin(ctx context.Context, auth *entity.IamUserAuth, identifier string, channel v1.LoginChannel, meta riskMeta) (*v1.TokenPair, *v1.SessionSummary, error) {
	_ = s.cache.ResetLoginFail(ctx, identifier)

	sid, err := s.newSID()
	if err != nil {
		return nil, nil, errs.Wrap(errs.CodeInternalError, err)
	}
	pair, refreshHash, _, refreshExp, err := s.issueTokenPair(auth.UserId, sid, uint32(auth.TokenVersion))
	if err != nil {
		return nil, nil, errs.Wrap(errs.CodeInternalError, err)
	}

	nowTime := gtime.NewFromTime(time.Now().UTC())
	err = dao.IamUserAuth.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.IamUserAuth.Columns()
		if _, err = tx.Model(dao.IamUserAuth.Table()).
			Where(cols.UserId, auth.UserId).
			Data(do.IamUserAuth{
				AccountStatus:        consts.AccountStatusActive,
				FailedLoginCount:     0,
				LastLoginAt:          nowTime,
				LastLoginIp:          meta.ClientIP,
				LastLoginGeo:         "",
				LastLoginUa:          meta.UserAgent,
				LastLoginFingerprint: meta.Fingerprint,
			}).
			Update(); err != nil {
			return err
		}

		_, err = tx.Model(dao.IamRefreshSession.Table()).Data(do.IamRefreshSession{
			Sid:              sid,
			UserId:           auth.UserId,
			RefreshTokenHash: refreshHash,
			UaHash:           sha256Hex(meta.UserAgent),
			Ip:               meta.ClientIP,
			Geo:              "",
			Fingerprint:      meta.Fingerprint,
			ExpiresAt:        gtime.NewFromTime(refreshExp.UTC()),
		}).Insert()
		return err
	})
	if err != nil {
		return nil, nil, errs.Wrap(errs.CodeInternalError, err)
	}

	authAfter, err := s.findAuthByUserID(ctx, auth.UserId)
	if err != nil || authAfter == nil {
		return nil, nil, errs.Wrap(errs.CodeInternalError, err)
	}
	session, err := s.buildSessionSummary(ctx, authAfter)
	if err != nil {
		return nil, nil, errs.Wrap(errs.CodeInternalError, err)
	}

	s.insertLoginLog(ctx, auth.UserId, identifier, channel, true, "", meta)
	return pair, session, nil
}
