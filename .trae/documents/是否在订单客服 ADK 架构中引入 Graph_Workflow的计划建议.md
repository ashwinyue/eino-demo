## 目标
- 在现有 ADK 架构基础上新增一个 AggregatorAgent，使用 Workflow 并行获取“订单详情、物流状态、政策说明”，统一格式化为单条回答，提高复杂询问的响应效率与质量。

## 改动范围
- 新增 `agents/aggregator.go`：实现 Workflow 并发与汇总，输出 `*schema.Message`。
- 更新 `agents/supervisor.go`：在路由提示中加入“当用户需要综合/汇总信息时分派给 aggregator”。
- 更新 `main.go`：注册 AggregatorAgent 到子 Agent 列表。

## 技术实现
- Workflow 结构：
  - 输入：`AggregatorInput{ OrderID string, Query string }`（从用户问题中解析出 `order_id`，或直接传入）
  - 节点：
    - `lambda_query_order`：调用 `tools.NewQueryOrderTool`（若缺参数调用澄清工具）
    - `lambda_track_delivery`：调用 `tools.NewTrackDeliveryTool`
    - `lambda_search_policy`：调用 `tools.NewSearchPolicyTool`
  - 汇总：`lambda_merge` 将三路结果格式化为 `*schema.Message`（有缺失项则用占位说明）。
- Workflow 编排：
  - `wf := compose.NewWorkflow[AggregatorInput, *schema.Message]()`
  - AddLambdaNode 三个并行节点；`wf.End()` 输入自 `lambda_merge`
- Agent 封装：
  - 使用 `adk.NewCustomAgent(ctx, &adk.CustomAgentConfig{...})` 或 `adk.NewChatModelAgent` 包装工作流执行逻辑，将用户输入转换为 `AggregatorInput` 并调用 `wfRunner.Invoke(ctx, input)`。

## 路由策略
- 在 `supervisor` 的指令中明确：
  - 当用户提“综合看订单情况/物流/政策建议/总体判断”等聚合需求时，分派到 `aggregator`。
  - 其他单一意图仍分派现有子 Agent（查询/取消/物流/发票/退款/售后/FAQ）。

## 验收
- 用示例问题验证聚合输出（例如“综合看下订单 20251112002 的状态、物流，并给我相关售后政策建议”）。
- 并发节点正常执行；缺失字段能触发澄清或给出占位提示；最终输出为一条结构化消息。

## 后续扩展
- 将政策检索替换为 RAG（Retriever + Indexer）；
- 引入 Interrupt/CheckPoint 以支持长流程暂停与恢复；
- 对并发节点增加失败重试与超时保护。