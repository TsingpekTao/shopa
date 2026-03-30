# Shopa Agent Execution Contract

本文件约束本仓库所有微服务开发流程。除非用户明确要求偏离，否则必须按以下顺序执行。

## 1. 统一开发流程
1. `gf init <service-name>` 初始化（新服务时）。
2. 先设计并落盘 `proto`（含 `go_package`、分页、幂等、状态字段语义）。
3. 执行 `gf gen pb -p manifest/protobuf/<domain> -a api -c internal/controller`。
4. 编写 `manifest/sql/0001_*.sql`（必要时追加 `0002_*.sql` 迁移）。
5. 执行 `gf gen dao` 生成 DAO/DO/Entity。
6. 在 `internal/logic` 实现业务逻辑（事务、幂等、状态机、补偿）。
7. 执行 `gf gen service` 生成并收口 service 接口。
8. Controller 对接 service，补齐路由与中间件。
9. 更新 `manifest/config/config.local|dev|prod.yaml`。
10. 冒烟验证 HTTP + gRPC + 跨服务链路。
11. 文档落盘到 `F:\shopa_doc\<service-name>\README.md` 与 `API_DETAIL.md`。

## 2. 代码结构约束
- Controller 只做协议层：参数绑定、鉴权上下文、错误映射。
- 业务逻辑必须在 `internal/logic`。
- `internal/logic/**` 新增/改动的有效语句行必须逐行中文注释（同一行尾注或上一行注释均可），并说明行为与业务意图/约束。
- 注释例外仅限：`import`、空行、仅花括号行、自动生成代码。
- 数据库写入只允许 `dao + do`，禁止 `g.Map`。
- 涉及库存/支付/退款的写链路必须具备：
  - 幂等键
  - 前置状态 CAS
  - 失败补偿路径

## 3. 命令规范
- Git Bash 启动：`GF_GCFG_FILE=manifest/config/config.local.yaml GFWORKER=off gf run main.go`
- PowerShell 启动：
  - `$env:GF_GCFG_FILE='manifest/config/config.local.yaml'`
  - `$env:GFWORKER='off'`
  - `gf run main.go`
- 构建验证：`go build ./...`
- 测试验证：`go test ./...`

## 4. Definition of Done
- `gf gen pb/dao/service` 全部成功。
- `go build ./...` 通过。
- 核心接口 curl 冒烟通过（至少 3 条：读、写、异常）。
- 对所有触达的 `internal/logic/**` 文件执行逐行中文注释审查，未满足规则视为不通过。
- 配置分环境可切换。
- 文档完成并与接口行为一致。

## 5. 安全与质量红线
- 严禁提交明文 AK/SK、DB 密码、私钥。
- 严禁把敏感配置写进 SQL。
- 严禁改动后不做编译验证。
- 严禁提交未满足逐行中文注释规范的 `internal/logic/**` 变更。
- 严禁为图省事绕过网关鉴权链路。

## 6. 前端协同要求
- 页面开发使用 `frontend-design` 风格规范。
- 提交前按 `web-design-guidelines` 做可用性与一致性检查。
- 前端对接接口前先冻结 API 契约（proto / swagger）。
