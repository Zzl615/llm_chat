# llm_chat

一个最小可运行（standalone）高并发 LLM 聊天 demo，使用 Go 语言实现。

## ✨ 特性

* **Gin + gorilla/websocket** 实现 RESTful & WebSocket 接入
* 每个会话（Session）拥有独立 goroutine，使用 channel 管理 I/O
* 模拟 **Kafka/Redis 队列**（用 Go channel 实现消息解耦）
* 流式返回（模拟大模型流式输出，定时推送生成内容）
* 支持 10 万连接的基础结构（单机 demo 可压测几千连接）
* 完整的单元测试覆盖

## 🧱 项目结构

```
llm_chat/
├── go.mod              # Go 模块定义
├── go.sum              # 依赖锁定文件
├── main.go             # 应用入口
├── internal/           # 内部包
│   ├── session.go      # Session 管理和 WebSocket 读写协程
│   ├── ws_handler.go   # Gin + WebSocket 路由处理器
│   └── mock_queue.go   # 模拟 Kafka/Redis 的消息队列
├── test/               # 单元测试
│   └── mock_queue_test.go  # MockQueue 测试用例
└── README.md           # 项目文档
```

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
- **功能**: 建立 WebSocket 连接，发送消息接收流式响应

### 消息格式

**客户端发送：**
```
纯文本消息（string）
```

**服务端响应：**
```
流式文本响应（string chunks）
格式：chunk 1: {消息内容}
      chunk 2: {消息内容}
      ...
      chunk 5: {消息内容}
```

## 🏗️ 架构设计

### 核心组件

1. **Session Manager**
   - 管理所有活跃的 WebSocket 会话
   - 每个会话拥有独立的读写协程
   - 使用 channel 实现异步消息传递

2. **MockQueue**
   - 模拟消息队列（类似 Kafka/Redis）
   - 实现请求-响应解耦
   - 支持异步处理和流式输出

3. **WebSocket Handler**
   - 基于 Gin 框架的路由处理
   - 自动处理连接升级
   - 支持并发连接管理

### 工作流程

```
客户端 → WebSocket 连接 → Session Manager
                              ↓
                         发布请求到队列
                              ↓
                       MockQueue Worker
                              ↓
                      模拟模型推理（流式输出）
                              ↓
                      结果投递到对应 Session
                              ↓
                      客户端接收流式响应
```

## 🔧 开发说明

### 代码风格

- 遵循 Go 官方代码规范
- 使用 `internal` 包隐藏内部实现
- 完整的注释和文档

### 依赖管理

使用 Go Modules，主要依赖：

- `github.com/gin-gonic/gin` - Web 框架
- `github.com/gorilla/websocket` - WebSocket 支持

## 📝 TODO

- [ ] 添加 RESTful API 端点
- [ ] 接入真实的Kafka队列
- [ ] 实现真实的 LLM 模型接入
- [ ] 添加连接限流和熔断机制
- [ ] 支持消息持久化
- [ ] 添加 Prometheus 监控指标
- [ ] 实现分布式部署支持

## 📄 License

MIT
