package auth

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/cache"
	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/mq"
	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	smsCodeMaxAttempts = 5
	smsCodeLength      = 6

	defaultAccessTTL  = 15 * time.Minute
	defaultRefreshTTL = 30 * 24 * time.Hour

	defaultLoginFailMax   = 5
	defaultLoginDelayBase = 200 * time.Millisecond
	defaultLoginDelayMax  = 2 * time.Second

	defaultUsernamePrefix = "\u7528\u6237_"

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

// jwtConf 保存 token 签发配置（密钥与过期时间）。
type jwtConf struct {
	Secret               string
	AccessExpireSeconds  int64
	RefreshExpireSeconds int64
}

// securityConf 保存登录风控开关与阈值配置。
type securityConf struct {
	LoginFailMax  int64
	MfaOnIPChange bool
}

// outboxConf 保存 outbox worker 的批量、退避、归档、清理参数。
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

// Service 聚合 IAM 业务依赖与运行配置，是 auth 领域的核心对象。
type Service struct {
	cache    *cache.Service
	mq       mq.Publisher
	node     *snowflake.Node
	jwt      jwtConf
	security securityConf
	outbox   outboxConf

	usernamePrefix string
	workerOnce     sync.Once
}

var (
	serviceOnce sync.Once
	serviceInst *Service
)

// New 返回 IAM 逻辑单例，保证后台 worker 和资源只初始化一次。
func New() *Service {
	serviceOnce.Do(func() {
		serviceInst = newService()
	})
	return serviceInst
}

// newService 从配置加载 IAM 运行参数并初始化依赖（cache、mq、snowflake）。
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

		node *snowflake.Node
		pub  mq.Publisher
		err  error
	)

	node, err = snowflake.NewNode(nodeID)
	if err != nil {
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
		g.Log().Warningf(ctx, "[iam-svc] init mq publisher failed, fallback noop: %+v", err)
		pub = mq.NewNoopPublisher()
	}

	return &Service{
		cache:          cache.New(),
		mq:             pub,
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

// refreshTTL 刷新凭证并返回新令牌。
func (s *Service) refreshTTL() time.Duration {
	return time.Duration(s.jwt.RefreshExpireSeconds) * time.Second
}

// cfgString 按 key 优先级读取字符串配置，读取不到则回落默认值。
func cfgString(ctx context.Context, def string, keys ...string) string {
	for _, key := range keys {
		v := strings.TrimSpace(g.Cfg().MustGet(ctx, key, "").String())
		if v != "" {
			return v
		}
	}
	return def
}

// cfgInt64 按 key 优先级读取 int64 配置，读取不到则回落默认值。
func cfgInt64(ctx context.Context, def int64, keys ...string) int64 {
	for _, key := range keys {
		v := g.Cfg().MustGet(ctx, key, int64(0)).Int64()
		if v != 0 {
			return v
		}
	}
	return def
}

// cfgBool 按 key 优先级读取布尔配置，读取不到则回落默认值。
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
