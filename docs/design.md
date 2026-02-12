# LLM Chat 架构设计文档

## 1. 项目简介

`llm_chat` 是一个基于 Go 语言的高并发 LLM 聊天演示项目。项目使用 **Gin + gorilla/websocket** 构建 WebSocket 服务，通过 Go channel 模拟消息队列（类 Kafka/Redis），实现大语言模型的流式响应效果。

本项目为纯内存架构，无数据库、无持久化存储、无真实 LLM 接入，所有推理过程均为模拟。

---

## 2. 技术栈

| 类别 | 技术选型 | 版本 |
|------|---------|------|
| 编程语言 | Go | 1.23.4 |
| Web 框架 | Gin | v1.10.0 |
| WebSocket | gorilla/websocket | v1.5.1 |
| 模块名称 | `llm-chat` | - |
| 构建工具 | Go Modules | - |

---

## 3. 目录结构

```
llm_chat/
├── main.go                    # 应用入口：初始化 Manager、MockQueue、路由，启动 HTTP 服务
├── go.mod                     # Go 模块定义及直接依赖声明
├── go.sum                     # 依赖版本锁定文件
├── internal/                  # 内部包（Go 强制私有，外部不可导入）
│   ├── manager.go             # Manager：线程安全的会话注册表
│   ├── session.go             # Session：单个 WebSocket 连接的读写协程管理
│   ├── mock_queue.go          # MockQueue：基于 channel 的模拟消息队列
│   └── ws_handler.go          # RegisterRoutes：Gin 路由注册及 WebSocket 升级
├── test/                      # 单元测试（独立包）
│   └── mock_queue_test.go     # MockQueue 的 5 个测试用例
├── docs/                      # 项目文档
│   └── design.md              # 架构设计文档（本文件）
├── CLAUDE.md                  # AI 辅助开发指导文件
├── README.md                  # 项目说明文档（中文）
└── .gitignore                 # Git 忽略规则
```

---

## 4. 核心架构

### 4.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                          客户端 (浏览器/wscat)                    │
└───────────────────────────────┬─────────────────────────────────┘
                                │ WebSocket (ws://localhost:8080/ws)
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Gin HTTP 服务器                          │
│                     ws_handler.go — 路由注册                     │
│            ┌──────────────────────────────────────┐              │
│            │  WebSocket Upgrader (协议升级)         │              │
│            │  CheckOrigin: 允许所有来源             │              │
│            └──────────────────┬───────────────────┘              │
└───────────────────────────────┼─────────────────────────────────┘
                                │ 创建 Session 并注册到 Manager
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Manager (会话管理器)                          │
│                      manager.go                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  sessions map[string]*Session    (sync.RWMutex 保护)     │    │
│  │  ├── sess-1 → *Session                                   │    │
│  │  ├── sess-2 → *Session                                   │    │
│  │  └── sess-N → *Session                                   │    │
│  └─────────────────────────────────────────────────────────┘    │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Session (会话实例)                           │
│                       session.go                                 │
│                                                                  │
│  ┌──────────────┐    Send chan []byte    ┌───────────────┐      │
│  │  ReadPump    │ ──────────────────────▶│  WritePump    │      │
│  │  (读协程)     │    (缓冲区: 128)        │  (写协程)      │      │
│  │              │                        │               │      │
│  │ 读取客户端消息 │                        │ 发送响应到客户端│      │
│  │ 转发到 Queue  │                        │ 定时发送 Ping  │      │
│  └──────┬───────┘                        └───────────────┘      │
└─────────┼───────────────────────────────────────────────────────┘
          │ handleRequest 回调
          ▼
┌─────────────────────────────────────────────────────────────────┐
│                  MockQueue (模拟消息队列)                         │
│                   mock_queue.go                                  │
│                                                                  │
│  ┌──────────────────┐         ┌──────────────────┐              │
│  │ requests channel  │ ──────▶│ MockModel Worker │              │
│  │ (缓冲区: 1024)    │         │ (模拟推理协程)    │              │
│  └──────────────────┘         └────────┬─────────┘              │
│                                        │ 每请求生成 5 个 chunk     │
│                                        │ 间隔 400ms              │
│                                        ▼                         │
│                               ┌──────────────────┐              │
│                               │ results channel   │              │
│                               │ (缓冲区: 1024)    │              │
│                               └────────┬─────────┘              │
│                                        │ SubscribeResults        │
│                                        ▼                         │
│                               回调 → Session.Send                │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 数据流

```
客户端发送消息
    │
    ▼
Session.ReadPump 读取 WebSocket 消息
    │
    ▼ handleRequest(sessionID, msg)
MockQueue.PublishRequest(Request{SessionID, Content})
    │
    ▼ requests channel
MockModel Worker 消费请求
    │
    ├── chunk 1 (400ms) ──▶ results channel
    ├── chunk 2 (400ms) ──▶ results channel
    ├── chunk 3 (400ms) ──▶ results channel
    ├── chunk 4 (400ms) ──▶ results channel
    └── chunk 5 (400ms) ──▶ results channel (IsLast=true)
                                │
                                ▼ SubscribeResults 回调
                          Manager.Get(sessionID)
                                │
                                ▼ 非阻塞写入
                          Session.Send channel
                                │
                                ▼
                          Session.WritePump 写入 WebSocket
                                │
                                ▼
                          客户端接收流式响应
```

---

## 5. 核心组件详解

### 5.1 Manager — 会话管理器

**文件**: `internal/manager.go` (40 行)

负责维护所有活跃 WebSocket 会话的注册表。

```go
type Manager struct {
    mu       sync.RWMutex
    sessions map[string]*Session
}
```

| 方法 | 锁类型 | 说明 |
|------|--------|------|
| `NewManager()` | - | 创建新的 Manager 实例 |
| `Register(s)` | 写锁 | 将 Session 注册到 sessions map |
| `Unregister(id)` | 写锁 | 从 sessions map 中移除指定 Session |
| `Get(id)` | 读锁 | 根据 ID 获取 Session，返回 `(*Session, bool)` |

**并发安全策略**: 使用 `sync.RWMutex` 实现读写分离，允许多个协程同时读取，写入时独占锁。

### 5.2 Session — 会话实例

**文件**: `internal/session.go` (93 行)

表示一个客户端 WebSocket 连接，每个 Session 启动两个 goroutine 分别处理读和写。

```go
type Session struct {
    ID      string
    Conn    *websocket.Conn
    Send    chan []byte     // 发送缓冲区，容量 128
    CloseCh chan struct{}   // 关闭信号通道
}
```

**时间常量**:

| 常量 | 值 | 说明 |
|------|----|------|
| `writeWait` | 5 秒 | 写操作超时时间 |
| `pongWait` | 30 秒 | 等待 Pong 响应超时 |
| `pingPeriod` | 27 秒 | Ping 发送间隔 (`pongWait * 9 / 10`) |

**ReadPump 工作流程**:

1. 设置读取限制为 1MB (`1 << 20`)
2. 设置读取超时为 `pongWait` (30 秒)
3. 注册 PongHandler 用于刷新读取超时
4. 循环读取消息，通过 `handleRequest` 回调转发到 MockQueue
5. 读取出错时关闭 `CloseCh` 和连接

**WritePump 工作流程**:

1. 创建 `pingPeriod` 定时器
2. 三路 `select` 监听:
   - `Send` channel：收到消息则写入 WebSocket
   - `ticker.C`：定时发送 Ping 帧保活
   - `CloseCh`：收到关闭信号则退出

### 5.3 MockQueue — 模拟消息队列

**文件**: `internal/mock_queue.go` (67 行)

用 Go channel 模拟 Kafka/Redis 风格的异步消息队列，实现请求与响应的解耦。

```go
type MockQueue struct {
    requests chan *Request   // 请求通道，缓冲区 1024
    results  chan *Result    // 结果通道，缓冲区 1024
}
```

**数据结构**:

```go
type Request struct {
    SessionID string    // 来源会话 ID
    Content   string    // 用户消息内容
}

type Result struct {
    SessionID string    // 目标会话 ID
    Chunk     string    // 流式输出片段
    IsLast    bool      // 是否为最后一个 chunk
}
```

| 方法 | 说明 |
|------|------|
| `NewMockQueue()` | 创建队列实例，初始化两个缓冲 channel |
| `StartMockModelWorker()` | 启动模型推理 worker 协程，消费 requests 生产 results |
| `PublishRequest(req)` | 将请求发布到 requests channel |
| `SubscribeResults(deliver)` | 启动订阅协程，将 results 逐个传递给回调函数 |

**模拟推理逻辑**: 每个请求生成 5 个 chunk，每个间隔 400ms，格式为 `chunk N: {原始内容}`。总响应时间约 2 秒。

### 5.4 WebSocket Handler — 路由处理器

**文件**: `internal/ws_handler.go` (37 行)

将 Gin HTTP 路由与 WebSocket 升级逻辑进行绑定。

**`RegisterRoutes` 处理流程**:

1. 注册 `GET /ws` 路由
2. 使用 `websocket.Upgrader` 将 HTTP 升级为 WebSocket
3. 生成 Session ID (`sess-{N}`)
4. 创建 Session 并注册到 Manager
5. 启动 WritePump 和 ReadPump goroutine

### 5.5 main.go — 应用入口

**文件**: `main.go` (39 行)

**启动流程**:

```
1. gin.Default()                    → 创建 Gin 引擎（含 Logger 和 Recovery 中间件）
2. NewManager()                     → 创建会话管理器
3. NewMockQueue()                   → 创建模拟消息队列
4. queue.StartMockModelWorker()     → 启动模型推理 worker
5. queue.SubscribeResults(callback) → 订阅结果，通过 Manager 分发到对应 Session
6. RegisterRoutes(r, manager, queue)→ 注册 WebSocket 路由
7. r.Run(":8080")                   → 启动 HTTP 服务，监听 8080 端口
```

**结果分发回调**:

```go
queue.SubscribeResults(func(res *internal.Result) {
    if sess, ok := manager.Get(res.SessionID); ok {
        select {
        case sess.Send <- []byte(res.Chunk):  // 非阻塞写入
        default:
            log.Printf("[WARN] session %s send buffer full, drop chunk", res.SessionID)
        }
    }
})
```

使用 `select + default` 实现非阻塞发送，当 Session 的 Send 缓冲区满时丢弃消息并打印警告日志，避免阻塞 SubscribeResults 协程。

---

## 6. 并发模型

### 6.1 Goroutine 分布

```
main goroutine
    │
    ├── MockQueue Worker goroutine      (1 个，处理所有请求)
    │
    ├── SubscribeResults goroutine       (1 个，分发所有结果)
    │
    └── 每个 WebSocket 连接
         ├── ReadPump goroutine          (1 个/连接)
         └── WritePump goroutine         (1 个/连接)
```

**总计**: N 个连接需要 `2N + 2` 个 goroutine（不含 Gin 框架自身的协程）。

### 6.2 Channel 通信矩阵

| Channel | 生产者 | 消费者 | 缓冲区大小 | 用途 |
|---------|--------|--------|-----------|------|
| `MockQueue.requests` | ReadPump (via PublishRequest) | MockModel Worker | 1024 | 请求排队 |
| `MockQueue.results` | MockModel Worker | SubscribeResults | 1024 | 结果分发 |
| `Session.Send` | SubscribeResults 回调 | WritePump | 128 | 单会话消息缓冲 |
| `Session.CloseCh` | ReadPump (close) | WritePump | 0 (无缓冲) | 关闭信号 |

### 6.3 共享状态保护

| 共享资源 | 保护机制 | 位置 |
|---------|---------|------|
| `Manager.sessions` | `sync.RWMutex` | manager.go |
| `Session.Conn` | 读写分离到不同 goroutine | session.go |
| channel 本身 | Go 运行时保证并发安全 | 全局 |

---

## 7. 协议与 API

### 7.1 WebSocket 端点

| 属性 | 值 |
|------|----|
| URL | `ws://localhost:8080/ws` |
| 协议 | RFC 6455 WebSocket |
| 方法 | GET (HTTP 升级) |
| CORS | 允许所有来源 |

### 7.2 消息格式

**客户端 → 服务端**:
```
纯文本字符串 (UTF-8)
```

**服务端 → 客户端** (流式响应):
```
chunk 1: {用户消息内容}
chunk 2: {用户消息内容}
chunk 3: {用户消息内容}
chunk 4: {用户消息内容}
chunk 5: {用户消息内容}
```

每个 chunk 间隔 400ms，共 5 个。最后一个 chunk 的 `Result.IsLast = true`。

### 7.3 WebSocket 保活机制

```
服务端 ──Ping──▶ 客户端    (每 27 秒)
服务端 ◀──Pong── 客户端    (自动响应)
```

如果 30 秒内未收到 Pong，ReadPump 读取超时，连接关闭。

---

## 8. 连接生命周期

```
1. 客户端发起 GET /ws 请求
       │
2. Gin 路由匹配，调用 handler
       │
3. WebSocket Upgrader 将 HTTP 升级为 WebSocket
       │
4. 生成 Session ID (sess-N)，创建 Session 实例
       │
5. Manager.Register(session) — 注册会话
       │
6. 启动 WritePump goroutine 和 ReadPump goroutine
       │
       ├── 正常阶段：ReadPump 读消息 → 发到队列 → 收到结果 → WritePump 写回
       │
7. 连接异常/客户端断开
       │
8. ReadPump 读取出错 → 关闭 CloseCh → 关闭 Conn
       │
9. WritePump 收到 CloseCh 信号 → 停止 ticker → 关闭 Conn
       │
      (注意: Manager.Unregister 当前未在断开时自动调用)
```

---

## 9. 构建与运行

```bash
# 安装依赖
go mod download

# 直接运行（开发模式）
go run .

# 编译二进制
go build -o llm_chat .

# 运行二进制
./llm_chat
```

服务启动后监听 `:8080` 端口。

**快速验证**:

```bash
# 使用 wscat
npm install -g wscat
wscat -c ws://localhost:8080/ws

# 或在浏览器控制台
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (e) => console.log(e.data);
ws.send('Hello, LLM!');
```

---

## 10. 测试

### 10.1 测试结构

测试位于 `test/` 目录下，使用独立的 `test` 包，仅依赖 Go 标准 `testing` 库。

```bash
# 运行所有测试
go test ./test -v

# 运行特定测试
go test ./test -v -run TestStartMockModelWorker

# 查看覆盖率
go test ./test -cover
```

### 10.2 测试用例

| 测试函数 | 验证内容 |
|---------|---------|
| `TestNewMockQueue` | MockQueue 初始化不为 nil |
| `TestPublishRequest` | PublishRequest 不阻塞 (1s 超时) |
| `TestStartMockModelWorker` | Worker 正确产生 5 个 chunk，SessionID 正确，IsLast 标志正确 |
| `TestMultipleRequests` | 3 个并发 Session 各自收到正确的 5 个 chunk |
| `TestConcurrentPublish` | 10 个并发 goroutine 同时发布请求不阻塞 |

### 10.3 异步测试模式

项目使用 `channel + select + time.After` 模式进行异步断言：

```go
select {
case <-done:
    // 验证结果
case <-time.After(5 * time.Second):
    t.Fatal("Timeout")
}
```

并发结果收集使用 `sync.Mutex` 保护共享的 map。

---

## 11. 代码规范

### 11.1 文件头

每个 Go 源文件顶部使用 JSDoc 风格注释：

```go
/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
```

### 11.2 命名规范

- **导出类型/函数**: PascalCase (`Manager`, `NewSession`, `RegisterRoutes`)
- **非导出字段/变量**: camelCase (`sessions`, `writeWait`, `requests`)
- **常量**: camelCase (`writeWait`, `pongWait`, `pingPeriod`)

### 11.3 日志格式

使用 `log.Printf` 配合方括号上下文标签：

```
[session sess-1] read err: ...
[MockModel] processing: sess-1 -> Hello
[WARN] session sess-1 send buffer full, drop chunk
```

### 11.4 包布局

- `internal/` — 所有非 main 代码，Go 强制外部不可导入
- `test/` — 独立包的单元测试，通过导出 API 进行黑盒测试

---

## 12. 已知限制

| 问题 | 位置 | 说明 |
|------|------|------|
| Session ID 竞态 | `ws_handler.go:28` | 读取 `mgr.sessions` map 长度时未持锁，并发场景下可能生成重复 ID |
| 无 CORS 限制 | `ws_handler.go:20` | `CheckOrigin` 对所有来源返回 `true` |
| 无优雅关闭 | `main.go` | 服务终止时不会等待活跃连接完成 |
| 串行处理 | `mock_queue.go` | MockModel Worker 为单协程，请求按 FIFO 顺序串行处理 |
| 会话泄露 | `ws_handler.go` | 连接断开时未调用 `Manager.Unregister`，session 残留在 map 中 |
| 无配置化 | `main.go` | 端口号、缓冲区大小等均为硬编码 |

---

## 13. 未来演进方向

根据 README 中的 TODO 规划：

| 方向 | 说明 |
|------|------|
| RESTful API | 增加 HTTP REST 端点 |
| Kafka 接入 | 替换 MockQueue 为真实 Kafka 队列 |
| LLM 接入 | 替换模拟推理为真实大语言模型调用 |
| 限流熔断 | 添加连接限流与熔断机制 |
| 消息持久化 | 引入数据库存储聊天记录 |
| 监控指标 | 接入 Prometheus 采集运行时指标 |
| 分布式部署 | 支持多实例水平扩展 |
