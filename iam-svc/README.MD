# iam-svc

IAM 认证微服务（GoFrame v2）。

## 配置切换

本服务已经将外部依赖配置集中放在 `manifest/config`，并提供多环境配置文件：

- `config.local.yaml`：本地开发
- `config.dev.yaml`：开发环境
- `config.prod.yaml`：生产环境
- `config.yaml`：默认入口（当前等价于 local）

外部组件相关配置包括：

- `database.default.link`
- `redis.default.*`
- `sms.*`
- `risk.*`
- `mq.rabbitmq.*`
- `geoip.*`
- `oauth.*`
- `iam.jwt.*` / `iam.security.*` / `iam.snowflake.*`

## 启动示例

### PowerShell（本地）

```powershell
cd F:\shopa\iam-svc
$env:GF_GCFG_FILE = "config.local.yaml"
gf run main.go
```

### 指定开发环境配置

```powershell
cd F:\shopa\iam-svc
$env:GF_GCFG_FILE = "config.dev.yaml"
gf run main.go
```

### Docker（生产配置）

容器默认 `GF_GCFG_FILE=config.prod.yaml`，也可以在运行时覆盖：

```bash
docker run -e GF_GCFG_FILE=config.prod.yaml <your-image>
```

## 说明

- IAM 服务代码优先读取 `iam.*` 命名空间配置，同时兼容旧键（`jwt.*`、`security.*`、`snowflake.*`）。
- 建议生产环境把敏感配置注入 Secret（K8s Secret / CI 变量 / 密钥管理服务），不要直接写入仓库。
