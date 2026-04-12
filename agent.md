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
- `internal/logic/**` 是业务逻辑强约束目录，新增或修改代码时必须执行“逐行中文注释”规范。
- 所谓“逐行中文注释”是指：每一条有效业务语句都必须有对应中文解释，允许写成上一行注释，或当前行尾注释。
- 注释内容不能只写“赋值”“调用方法”这类空话，必须明确说明：
  - 这行代码在做什么
  - 为什么要这么做
  - 它对应的业务规则、状态约束、幂等/补偿/风控语义是什么
- 注释例外仅限：
  - `import`
  - 空行
  - 仅花括号行
  - 结构体字段标签
  - 接口/类型定义中无业务语义的纯声明行
  - 自动生成代码
- 如果某段逻辑是多行联动，仍然不能只在代码块开头写一条笼统注释，必须让阅读者能逐行看懂每一步业务意图。
- 对 `internal/logic/**` 的 CR/Review，逐行中文注释属于阻塞项；未满足时视为代码不合格，不能提交、不能合并。
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

## 3.1 VS Code MCP / 集成终端约束
- 当 Codex 已接入可操作 VS Code 集成终端的 MCP server 时，凡是“启动运行态进程”的动作，必须优先在 VS Code 界面的终端中执行，不得默认改在 Codex 自己的 shell 中启动。
- “启动运行态进程”仅限以下场景：
  - 启动任一微服务（如 `gf run main.go`、`go run main.go`、项目约定的服务启动脚本）。
  - 启动前端开发服务（如 `npm run dev`、`pnpm dev`、`yarn dev`、项目约定的 web 启动脚本）。
- VS Code 终端用途严格限制为“启动微服务”和“启动前端服务”，不得把下列操作挪到 VS Code 终端执行：
  - 代码搜索、文件检查、脚本探查。
  - `go build`、`go test`、`gf gen`、`npm install`、`pnpm install` 等构建/生成/安装命令。
  - `git status`、`git diff`、`git add`、`git commit` 等版本控制命令。
- 除“启动微服务/前端服务”外，其余命令默认仍在 Codex 当前工作环境的 shell 中执行。
- 如果当前会话尚未接入 VS Code 终端 MCP server，必须先明确告知用户“暂时无法把启动命令投递到 VS Code 终端”，不得假装已在 VS Code 中启动。
- 在未接入 VS Code 终端 MCP server 的前提下，除非用户当次明确允许回退到 Codex shell 启动服务，否则不要擅自改用当前 shell 挂起长期运行进程。
- 若用户要求启动多个运行态进程，应分别在 VS Code 中使用独立终端承载，避免把多个服务复用到同一个终端会话里。

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
- 面向用户的前端文案必须使用“产品话术”，禁止出现任何面向程序员/开发调试的提示。
- 严禁在页面、弹窗、空态、按钮说明、帮助文案中出现以下类型的话术：
  - “当前版本”“V1/一期/二期”“开发中”“联调中”“待接入”“暂未对接真实能力”
  - “配置 baseUrl / API Key / token / mock 数据后可用”
  - “接口返回”“字段缺失”“轮询中”“降级中”“占位数据”“测试环境”
  - 任何把实现状态、技术方案、配置要求直接暴露给用户的程序员视角表述
- 如果真实能力尚未接通，前端也必须改写成用户可理解的话术，例如“当前服务暂不可用，请稍后再试”或“正在为你连接客服，请稍候”，禁止把内部原因直接展示给用户。
- Code Review 时，凡发现前端文案含开发者提示、调试说明、技术配置词汇，视为阻塞项，必须修改后才能提交或合并。
