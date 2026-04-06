package auth

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/cache"
	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/mq"
	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/sms"
	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	// 短信验证码相关安全参数。
	smsCodeMaxAttempts = 5
	smsCodeLength      = 6

	// JWT 默认有效期策略：access 较短、refresh 较长。
	defaultAccessTTL  = 15 * time.Minute
	defaultRefreshTTL = 30 * 24 * time.Hour

	// 登录失败控制参数，用于锁定和退避。
	defaultLoginFailMax   = 5
	defaultLoginDelayBase = 200 * time.Millisecond
	defaultLoginDelayMax  = 2 * time.Second

	// 新注册用户的默认展示名前缀。
	defaultUsernamePrefix = "用户_"

	// Outbox worker 默认运行参数。
	defaultOutboxBatchSize        = 100
	defaultOutboxPollInterval     = 1 * time.Second
	defaultOutboxMaxFailCount     = 20
	defaultOutboxFailBaseDelay    = 5 * time.Second
	defaultOutboxFailMaxDelay     = 30 * time.Minute
	defaultOutboxArchiveBatchSize = 1000
	defaultOutboxArchiveAfter     = 7 * 24 * time.Hour
	defaultOutboxArchiveInterval  = 5 * time.Minute
	defaultOutboxCleanupBatchSize = 1000
	defaultOutboxCleanupAfter     = 90 * 24 * time.Hour
	defaultOutboxCleanupInterval  = 30 * time.Minute
	defaultOutboxMetricsInterval  = 1 * time.Minute
)

// jwtConf 保存 token 签发配置。
type jwtConf struct {
	Secret               string
	AccessExpireSeconds  int64
	RefreshExpireSeconds int64
}

// securityConf 保存登录安全策略配置。
type securityConf struct {
	LoginFailMax  int64
	MfaOnIPChange bool
}

// outboxConf 保存 outbox 后台任务运行参数。
type outboxConf struct {
	BatchSize       int
	PollInterval    time.Duration
	MaxFailCount    int
	FailBaseDelay   time.Duration
	FailMaxDelay    time.Duration
	ArchiveBatch    int
	ArchiveAfter    time.Duration
	ArchiveInterval time.Duration
	CleanupBatch    int
	CleanupAfter    time.Duration
	CleanupInterval time.Duration
	MetricsInterval time.Duration
}

// Service 是 auth 领域的核心服务，聚合运行依赖和配置。
// workerOnce 用于保证后台 worker 只启动一次，避免重复消费 outbox。
type Service struct {
	cache    *cache.Service
	mq       mq.Publisher
	sms      sms.Sender
	node     *snowflake.Node
	jwt      jwtConf
	security securityConf
	outbox   outboxConf

	usernamePrefix string
	workerOnce     sync.Once
}

var (
	// serviceOnce 保证 Service 单例仅初始化一次。
	serviceOnce sync.Once
	serviceInst *Service
)

// New 返回 IAM 逻辑单例。
func New() *Service {
	serviceOnce.Do(func() {
		serviceInst = newService()
	})
	return serviceInst
}

// newService 从配置加载运行参数，并初始化 cache、mq、snowflake 等依赖。
func newService() *Service {
	var (
		ctx = context.Background()

		jwtCfg = jwtConf{
			Secret: cfgString(ctx, "", "iam.jwt.secret", "jwt.secret"),
			AccessExpireSeconds: cfgInt64(
				ctx,
				int64(defaultAccessTTL/time.Second),
				"iam.jwt.accessExpireSeconds",
				"jwt.accessExpireSeconds",
			),
			RefreshExpireSeconds: cfgInt64(
				ctx,
				int64(defaultRefreshTTL/time.Second),
				"iam.jwt.refreshExpireSeconds",
				"jwt.refreshExpireSeconds",
			),
		}
		secCfg = securityConf{
			LoginFailMax: cfgInt64(ctx, defaultLoginFailMax, "iam.security.loginFailMax", "security.loginFailMax"),
			MfaOnIPChange: cfgBool(
				ctx,
				true,
				"iam.security.mfaOnIpChange",
				"security.mfaOnIpChange",
			),
		}
		outCfg = outboxConf{
			BatchSize:       int(cfgInt64(ctx, defaultOutboxBatchSize, "iam.outbox.batchSize")),
			PollInterval:    time.Duration(cfgInt64(ctx, int64(defaultOutboxPollInterval/time.Millisecond), "iam.outbox.pollIntervalMs")) * time.Millisecond,
			MaxFailCount:    int(cfgInt64(ctx, defaultOutboxMaxFailCount, "iam.outbox.maxFailCount")),
			FailBaseDelay:   time.Duration(cfgInt64(ctx, int64(defaultOutboxFailBaseDelay/time.Millisecond), "iam.outbox.failBaseDelayMs")) * time.Millisecond,
			FailMaxDelay:    time.Duration(cfgInt64(ctx, int64(defaultOutboxFailMaxDelay/time.Millisecond), "iam.outbox.failMaxDelayMs")) * time.Millisecond,
			ArchiveBatch:    int(cfgInt64(ctx, defaultOutboxArchiveBatchSize, "iam.outbox.archiveBatchSize")),
			ArchiveAfter:    time.Duration(cfgInt64(ctx, int64(defaultOutboxArchiveAfter/time.Hour), "iam.outbox.archiveAfterHours")) * time.Hour,
			ArchiveInterval: time.Duration(cfgInt64(ctx, int64(defaultOutboxArchiveInterval/time.Second), "iam.outbox.archiveIntervalSeconds")) * time.Second,
			CleanupBatch:    int(cfgInt64(ctx, defaultOutboxCleanupBatchSize, "iam.outbox.cleanupBatchSize")),
			CleanupAfter:    time.Duration(cfgInt64(ctx, int64(defaultOutboxCleanupAfter/time.Hour), "iam.outbox.cleanupAfterHours")) * time.Hour,
			CleanupInterval: time.Duration(cfgInt64(ctx, int64(defaultOutboxCleanupInterval/time.Second), "iam.outbox.cleanupIntervalSeconds")) * time.Second,
			MetricsInterval: time.Duration(cfgInt64(ctx, int64(defaultOutboxMetricsInterval/time.Second), "iam.outbox.metricsIntervalSeconds")) * time.Second,
		}

		nodeID     = cfgInt64(ctx, 1, "iam.snowflake.node", "snowflake.node")
		namePrefix = cfgString(ctx, defaultUsernamePrefix, "iam.register.defaultUsernamePrefix", "register.defaultUsernamePrefix")

		node      *snowflake.Node
		pub       mq.Publisher
		smsSender sms.Sender
		err       error
	)

	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		// 节点号非法时回退到 1，避免服务因配置问题完全不可用。
		node, _ = snowflake.NewNode(1)
	}
	if jwtCfg.AccessExpireSeconds <= 0 {
		jwtCfg.AccessExpireSeconds = int64(defaultAccessTTL / time.Second)
	}
	if jwtCfg.RefreshExpireSeconds <= 0 {
		jwtCfg.RefreshExpireSeconds = int64(defaultRefreshTTL / time.Second)
	}
	if secCfg.LoginFailMax <= 0 {
		secCfg.LoginFailMax = defaultLoginFailMax
	}
	if jwtCfg.Secret == "" {
		// 本地开发环境允许使用兜底 secret，避免因缺配置导致服务无法启动。
		jwtCfg.Secret = "shopa-iam-dev-secret"
	}

	if outCfg.BatchSize <= 0 {
		outCfg.BatchSize = defaultOutboxBatchSize
	}
	if outCfg.PollInterval <= 0 {
		outCfg.PollInterval = defaultOutboxPollInterval
	}
	if outCfg.MaxFailCount <= 0 {
		outCfg.MaxFailCount = defaultOutboxMaxFailCount
	}
	if outCfg.FailBaseDelay <= 0 {
		outCfg.FailBaseDelay = defaultOutboxFailBaseDelay
	}
	if outCfg.FailMaxDelay <= 0 {
		outCfg.FailMaxDelay = defaultOutboxFailMaxDelay
	}
	if outCfg.ArchiveBatch <= 0 {
		outCfg.ArchiveBatch = defaultOutboxArchiveBatchSize
	}
	if outCfg.ArchiveAfter <= 0 {
		outCfg.ArchiveAfter = defaultOutboxArchiveAfter
	}
	if outCfg.ArchiveInterval <= 0 {
		outCfg.ArchiveInterval = defaultOutboxArchiveInterval
	}
	if outCfg.CleanupBatch <= 0 {
		outCfg.CleanupBatch = defaultOutboxCleanupBatchSize
	}
	if outCfg.CleanupAfter <= 0 {
		outCfg.CleanupAfter = defaultOutboxCleanupAfter
	}
	if outCfg.CleanupInterval <= 0 {
		outCfg.CleanupInterval = defaultOutboxCleanupInterval
	}
	if outCfg.MetricsInterval <= 0 {
		outCfg.MetricsInterval = defaultOutboxMetricsInterval
	}

	pub, err = mq.NewPublisherFromConfig(ctx)
	if err != nil {
		// MQ 初始化失败时退化为 noop publisher，保证主链路仍可运行。
		g.Log().Warningf(ctx, "[iam-svc] init mq publisher failed, fallback noop: %+v", err)
		pub = mq.NewNoopPublisher()
	}
	smsSender, err = sms.NewSenderFromConfig(ctx)
	if err != nil {
		// 短信渠道初始化失败时回退为 mock sender，方便本地联调。
		g.Log().Warningf(ctx, "[iam-svc] init sms sender failed, fallback mock: %+v", err)
		smsSender = sms.NewMockSender()
	}

	return &Service{
		cache:          cache.New(),
		mq:             pub,
		sms:            smsSender,
		node:           node,
		jwt:            jwtCfg,
		security:       secCfg,
		outbox:         outCfg,
		usernamePrefix: namePrefix,
	}
}

// accessTTL 返回 access token 的有效期。
func (s *Service) accessTTL() time.Duration {
	return time.Duration(s.jwt.AccessExpireSeconds) * time.Second
}

// refreshTTL 返回 refresh token 的有效期。
func (s *Service) refreshTTL() time.Duration {
	return time.Duration(s.jwt.RefreshExpireSeconds) * time.Second
}

// cfgString 依次读取多个 key，返回第一个非空字符串配置。
func cfgString(ctx context.Context, def string, keys ...string) string {
	for _, key := range keys {
		v := strings.TrimSpace(g.Cfg().MustGet(ctx, key, "").String())
		if v != "" {
			return v
		}
	}
	return def
}

// cfgInt64 依次读取多个 key，返回第一个非零 int64 配置。
func cfgInt64(ctx context.Context, def int64, keys ...string) int64 {
	for _, key := range keys {
		v := g.Cfg().MustGet(ctx, key, int64(0)).Int64()
		if v != 0 {
			return v
		}
	}
	return def
}

// cfgBool 依次读取多个 key，只要显式配置了布尔值就返回。
func cfgBool(ctx context.Context, def bool, keys ...string) bool {
	for _, key := range keys {
		v := g.Cfg().MustGet(ctx, key, nil)
		if v.IsNil() {
			continue
		}
		raw := strings.TrimSpace(v.String())
		if raw == "" {
			continue
		}
		return v.Bool()
	}
	return def
}
