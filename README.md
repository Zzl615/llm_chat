# llm_chat

一个最小可运行（standalone）高并发 LLM 聊天 demo，使用 Go 语言实现，采用 **DDD（领域驱动设计）** 架构。

## ✨ 特性

* **DDD 四层架构**：Domain、Application、Infrastructure、Interface 清晰分层
* **Gin + gorilla/websocket** 实现 RESTful & WebSocket 接入
* 每个会话（Session）拥有独立 goroutine，使用 channel 管理 I/O
* 模拟 **Kafka/Redis 队列**（用 Go channel 实现消息解耦）
* 流式返回（模拟大模型流式输出，定时推送生成内容）
* 支持 10 万连接的基础结构（单机 demo 可压测几千连接）
* 完整的单元测试覆盖

## 🧱 项目结构（DDD 架构）

```
llm-chat/
├── go.mod              # Go 模块定义
├── go.sum              # 依赖锁定文件
├── main.go             # 应用入口（依赖注入和组装）
├── domain/             # 领域层（核心业务逻辑）
│   ├── entity/         # 领域实体
│   │   ├── session.go  # Session 实体
│   │   └── message.go  # Message 实体
│   ├── valueobject/    # 值对象
│   │   └── session_id.go
│   ├── service/        # 领域服务
│   │   └── session_service.go
│   └── repository/     # 仓储接口（依赖倒置）
│       └── session_repository.go
├── application/        # 应用层（用例编排）
│   ├── service/        # 应用服务
│   │   ├── chat_service.go
│   │   └── session_service.go
│   └── dto/            # 数据传输对象
│       ├── request.go
│       └── response.go
├── infrastructure/     # 基础设施层（技术实现）
│   ├── repository/     # 仓储实现
│   │   └── session_repository_impl.go
│   ├── queue/          # 消息队列实现
│   │   └── mock_queue.go
│   └── websocket/      # WebSocket 基础设施
│       └── connection.go
├── interface/          # 接口层（外部接口）
│   ├── http/           # HTTP 路由
│   │   └── router.go
│   ├── websocket/      # WebSocket 处理器
│   │   └── handler.go
│   └── sse/            # SSE 处理器
│       └── handler.go
├── test/               # 单元测试
│   └── mock_queue_test.go
└── README.md           # 项目文档
```

## 🏗️ DDD 架构说明

### 依赖方向

```
Interface → Application → Domain
                ↑
         Infrastructure ──┘
```

- **Domain（领域层）**：核心业务逻辑，不依赖任何外部技术
- **Application（应用层）**：用例编排，依赖 Domain
- **Infrastructure（基础设施层）**：技术实现，实现 Domain 和 Application 定义的接口
- **Interface（接口层）**：外部接口（HTTP、WebSocket、SSE），依赖 Application

### 各层职责

1. **Domain Layer（领域层）**
   - `entity/`：领域实体（Session、Message）
   - `valueobject/`：值对象（SessionID）
   - `service/`：领域服务（SessionService）
   - `repository/`：仓储接口（依赖倒置）

2. **Application Layer（应用层）**
   - `service/`：应用服务（ChatService、SessionApplicationService）
   - `dto/`：数据传输对象（Request、Response）

3. **Infrastructure Layer（基础设施层）**
   - `repository/`：仓储实现（SessionRepositoryImpl）
   - `queue/`：消息队列实现（MockQueue）
   - `websocket/`：WebSocket 连接封装

4. **Interface Layer（接口层）**
   - `http/`：HTTP 路由注册
   - `websocket/`：WebSocket 处理器
   - `sse/`：SSE 处理器

## 🚀 快速开始

### 安装依赖

```bash
go mod download
```

### 运行服务

```bash
go run .
```

服务将在 `:8080` 端口启动。

### 测试 WebSocket 连接

访问：[http://localhost:8080/ws](http://localhost:8080/ws)

**使用 WebSocket 客户端测试：**

```javascript
// 浏览器 Console
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
    console.log('Received:', event.data);
};
ws.send('Hello, LLM!');
```

**使用 wscat 命令行工具：**

```bash
# 安装 wscat
npm install -g wscat

# 连接测试
wscat -c ws://localhost:8080/ws
```

### 测试 SSE 连接

访问：[http://localhost:8080/sse/{sessionId}](http://localhost:8080/sse/{sessionId})

**使用 curl 测试：**

```bash
curl -N http://localhost:8080/sse/sess-1
```

## 🧪 运行测试

```bash
# 运行所有测试
go test ./test -v

# 运行特定测试
go test ./test -v -run TestNewMockQueue

# 查看测试覆盖率
go test ./test -cover
```

## 📖 API 说明

### WebSocket 端点

- **URL**: `/ws`
- **协议**: WebSocket
- **功能**: 建立 WebSocket 连接，创建会话并发送消息

### SSE 端点

- **URL**: `/sse/:sessionId`
- **协议**: Server-Sent Events
- **功能**: 接收大模型返回的流式结果

### 消息格式

**客户端发送（WebSocket）：**
```
纯文本消息（string）
```

**服务端响应（SSE）：**
```
流式文本响应（JSON chunks）
格式：{"chunk": "chunk 1: {消息内容}", "is_last": false}
      {"chunk": "chunk 2: {消息内容}", "is_last": false}
      ...
      {"chunk": "chunk 5: {消息内容}", "is_last": true}
```

## 🏗️ 架构设计

### 核心组件

1. **Session Entity（会话实体）**
   - 领域实体，包含业务属性
   - 不包含技术细节（如 WebSocket 连接）

2. **SessionRepository（会话仓储）**
   - 仓储接口定义在 Domain 层
   - 仓储实现放在 Infrastructure 层
   - 实现依赖倒置原则

3. **MessageQueue（消息队列）**
   - 消息队列接口定义在 Application 层
   - MockQueue 实现放在 Infrastructure 层
   - 实现请求-响应解耦

4. **WebSocket/SSE Handler（处理器）**
   - 接口层处理外部连接
   - 依赖注入应用服务

### 工作流程

```
客户端 → WebSocket 连接 → Interface Layer
                              ↓
                      Application Layer (ChatService)
                              ↓
                      Domain Layer (Session Entity)
                              ↓
                      Infrastructure Layer (Repository)
                              ↓
                      Application Layer (MessageQueue)
                              ↓
                      Infrastructure Layer (MockQueue)
                              ↓
                      模拟模型推理（流式输出）
                              ↓
                      SSE 处理器订阅结果
                              ↓
                      客户端接收流式响应
```

## 🔧 开发说明

### 代码风格

- 遵循 Go 官方代码规范
- 采用 DDD 四层架构，职责清晰
- 完整的注释和文档
- 依赖倒置原则（Domain 定义接口，Infrastructure 实现）

### 依赖管理

使用 Go Modules，主要依赖：

- `github.com/gin-gonic/gin` - Web 框架
- `github.com/gorilla/websocket` - WebSocket 支持

### DDD 设计原则

1. **依赖倒置**：Domain 层定义接口，Infrastructure 层实现
2. **领域驱动**：核心业务逻辑集中在 Domain 层
3. **分层清晰**：各层职责明确，依赖方向单一
4. **接口隔离**：使用接口实现解耦

## 📝 TODO

- [ ] 添加 RESTful API 端点
- [ ] 接入真实的 Kafka 队列
- [ ] 实现真实的 LLM 模型接入
- [ ] 添加连接限流和熔断机制
- [ ] 支持消息持久化
- [ ] 添加 Prometheus 监控指标
- [ ] 实现分布式部署支持
- [ ] 添加领域事件（Domain Events）

## 📄 License

MIT
