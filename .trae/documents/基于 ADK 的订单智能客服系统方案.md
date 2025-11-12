## 目标
- 使用 Eino 的 ADK 模式构建完善的多 Agent 订单智能客服系统，覆盖查询、取消、售后、物流、发票、支付/退款与升级处理。
- 根据场景灵活选择 ADK 的主管路由、顺序/并行/循环工作流与 Plan-Execute-Replan 策略，支持流式输出、工具调用与会话记忆。

## 总体架构
- 顶层主管（Supervisor Agent）：根据用户意图路由到合适子 Agent，避免自行执行任务，参考 `adk/multiagent/supervisor`。
- 子 Agent（ChatModelAgent + Tools）：每个域（订单查询、取消、售后、FAQ、物流、发票、支付/退款、升级）独立实现，绑定对应工具。
- 数据与工具层：通过 `utils.InferOptionableTool` 封装订单服务/物流/发票等接口；提供 UnknownToolsHandler 与错误处理。
- 工作流层：在复杂会话中使用 ADK 工作流（顺序/并行/循环）组织多步骤任务；对需要字段映射的场景使用 Workflow；对需要中断与恢复的任务启用 CheckPoint。
- 会话与观察：启用会话状态管理与链路追踪（参考 `adk/common/session`、`adk/common/trace`），支持流式查询与事件打印（参考 `adk/intro/workflow/sequential` 的 Runner 用法）。

## Agent 划分
- OrderQueryAgent：查询订单详情、状态与商品列表
- OrderCancelAgent：校验状态并执行取消
- AfterSalesAgent：基于知识库/RAG提供退换货与质保政策
- FAQAgent：常见问题（发票、配送、时效、售后入口等）
- LogisticsAgent：查询物流状态、预计送达
- InvoiceAgent：开具/补开发票规则说明或触发工具
- PaymentRefundAgent：支付状态与退款申请流程（含必要校验）
- EscalationAgent：无法自动解决时，输出升级路径或回传人工坐席入口

## 工具设计（示例）
- `query_order`：输入 `order_id`，返回订单概要；绑定到 OrderQueryAgent
- `cancel_order`：输入 `order_id`，执行取消并返回结果；绑定到 OrderCancelAgent
- `track_delivery`：输入 `order_id` 或 `tracking_no`，返回物流信息；绑定到 LogisticsAgent
- `issue_invoice`：输入 `order_id`、`title`、`tax_id`，返回开票结果/链接；绑定到 InvoiceAgent
- `apply_refund`：输入 `order_id`、`reason`，返回受理结果；绑定到 PaymentRefundAgent
- `ask_for_clarification`：缺失关键参数时追问用户，参考 `plan-execute-replan/tools/ask_for_clarification.go`
- 工具封装统一采用 `utils.InferOptionableTool`，并定义清晰的 `json` schema 字段与错误语义。

## 意图识别与路由
- 轻量策略：规则分类（关键词/正则）作为兜底路径，避免 LLM 不可用时阻塞
- LLM 策略：主管 Agent 用提示工程明确“只路由不执行”，子 Agent 各自执行工具；必要时引入单独的分类器 Agent
- 分派策略：一次只分派一个 Agent，避免并行工具滥用；需要并发时转入并行工作流层处理

## 工作流模式
- 顺序（Sequential）：多步骤串行（例如“补全缺失参数→查询→格式化响应”），参考 `intro/workflow/sequential`
- 并行（Parallel）：同时拉取订单信息与物流状态，参考 `intro/workflow/parallel`
- 循环（Loop/Reflection）：无法回答时反思与追问，参考 `intro/workflow/loop`
- Plan-Execute-Replan：复杂任务由 Planner 生成计划，Executor 执行，Replanner 根据结果调整，参考 `multiagent/plan-execute-replan`
- 中断与检查点：对长流程开启 Interrupt/CheckPoint，支持暂停与恢复（参考文档新增的 CheckPoint/Interrupt 指南）

## 数据与知识库
- 订单服务：对接内部订单 API（REST/gRPC），工具中处理鉴权、重试与限流
- RAG：将售后政策、FAQ、发票流程等文档建立索引（ES/VikingDB），Retriever 按查询返回上下文
- 文档加载：用 `Document Loader` 将政策/帮助文档加载到索引；Transformer 做清洗与分段

## 会话与状态
- 会话上下文保存最近 N 轮消息、用户标识、已选订单号等关键状态
- 在工具或工作流节点间传递结构化状态（如 `RouteInfo`、`OrderContext`）
- 必要时对状态进行序列化并持久化（CheckPoint）

## 观测与运维
- 事件打印与链路追踪：沿用 `adk/common/trace` 的 Coze Loop 回调与 `common/prints`
- 日志与指标：为工具调用、错误、重试、耗时暴露日志与指标；记录路由决策与执行结果
- 安全与合规：隐藏 PII，校验订单号格式，避免在日志中记录敏感信息；对外调用加鉴权

## 配置与环境
- 模型：通过环境选择 OpenAI 或 Ark（参考 `adk/common/model/chat_model.go`）
- 服务端点：订单、物流、发票、支付等服务地址与密钥通过环境变量配置
- Runner：默认开启流式输出（可配置），支持超时、重试策略

## 测试计划
- 单测：工具函数（入参校验、错误处理、重试策略）、意图分类与路由
- 集成测：对接模拟订单/物流服务，验证各 Agent 的工具调用路径
- 端到端：使用 Runner 执行真实对话样例，覆盖查询/取消/售后/物流/发票/退款/升级

## 目录结构（对齐示例风格）
- `adk/orders-cs/`
  - `agents/`：各子 Agent 与主管 Agent 定义
  - `tools/`：订单、取消、物流、发票、退款、澄清等工具
  - `workflows/`：顺序/并行/循环与复合编排
  - `common/`：模型选择、会话状态、打印与追踪
  - `main.go`：集成 Runner 入口（支持流式打印）

## 里程碑
1. 核心功能（主管路由 + 查询/取消）：可运行，工具接入模拟数据源
2. 售后/FAQ + RAG：接入政策文档并检索回答
3. 物流/发票/支付退款：补齐工具与子 Agent
4. 复杂任务（Plan-Execute-Replan）与 Clarification 工具
5. 会话管理与 CheckPoint/Interrupt 支持
6. 接入真实服务端点与鉴权、安全加固、观测完善

## 验收标准
- 覆盖主要订单需求场景并通过端到端测试
- 路由准确率与错误处理健壮；对无关键参数能主动澄清
- 性能与稳定性满足并发需求；日志与追踪可定位问题

## 下一步
- 按上面目录创建模块骨架与最小可运行样例（主管路由 + 查询/取消）
- 在本仓库以示例形式落地，后续迭代补齐 RAG、工作流与复杂策略