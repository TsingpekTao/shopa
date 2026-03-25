# 配置文件说明

## 文件列表

- `config.yaml`：默认配置入口（未指定 `GF_GCFG_FILE` 时使用）。
- `config.local.yaml`：本地开发配置。
- `config.dev.yaml`：开发环境配置。
- `config.prod.yaml`：生产环境配置。

## 切换方式

GoFrame 支持通过 `GF_GCFG_FILE` 指定配置文件名：

```powershell
$env:GF_GCFG_FILE = "config.dev.yaml"
```

或命令行参数：

```bash
./main --gf.gcfg.file=config.prod.yaml
```

## 推荐实践

- 数据库、Redis、短信、风控、MQ、OAuth 等第三方配置统一维护在该目录。
- 生产环境敏感信息优先通过 Secret 注入，再由部署流程渲染到配置文件。
- `iam-svc` 业务配置优先使用 `iam.*` 命名空间，避免后续多服务共享配置时冲突。
