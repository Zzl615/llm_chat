# llm_chat

一个基于 Kratos 框架的高并发 LLM 聊天服务，使用 Go 语言实现。

## ✨ 特性

* **Kratos 微服务框架** - 使用 Kratos v2 框架，遵循标准项目结构
* **WebSocket 支持** - 基于 gorilla/websocket 实现实时通信
* **分层架构** - 采用 DDD 分层架构（biz/service/server）
* **消息队列** - 模拟 Kafka/Redis 队列实现消息解耦
* **流式返回** - 模拟大模型流式输出，定时推送生成内容
* **高并发支持** - 支持大量并发连接的基础结构
* **配置管理** - 使用 YAML 配置文件，支持动态配置
* **优雅关闭** - 支持优雅关闭和资源清理

## 🧱 项目结构

```
llm_chat/
├── cmd/                    # 应用入口
│   └── server/             # 服务器入口
│       └── main.go
├── configs/                # 配置文件
│   └── config.yaml         # 应用配置
├── internal/               # 内部代码
│   ├── biz/               # 业务逻辑层
│   │   ├── manager.go     # Session 管理器
│   │   ├── queue.go       # 消息队列接口和实现
│   │   └── session.go     # WebSocket 会话管理
│   ├── conf/              # 配置定义
│   │   ├── config.go      # 配置结构体
│   │   └── config.proto    # 配置 Proto 定义（可选）
│   ├── server/            # 传输层（HTTP/gRPC）
│   │   └── http.go        # HTTP 服务器
│   └── service/           # 服务层
│       └── chat.go        # 聊天服务
├── test/                   # 单元测试
│   └── mock_queue_test.go # MockQueue 测试用例
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖锁定文件
└── README.md               # 项目文档
```

## 🚀 快速开始

### 前置要求

- Go 1.23+ 
- 已安装 Go 模块支持

### 安装依赖

```bash
go mod download
```

### 运行服务

```bash
# 从项目根目录运行
go run cmd/server/main.go

# 或指定配置文件路径
go run cmd/server/main.go -conf configs/config.yaml
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

1. **Session Manager (biz.Manager)**
   - 管理所有活跃的 WebSocket 会话
   - 每个会话拥有独立的读写协程
   - 使用 channel 实现异步消息传递

2. **MockQueue (biz.MockQueue)**
   - 模拟消息队列（类似 Kafka/Redis）
   - 实现请求-响应解耦
   - 支持异步处理和流式输出

3. **ChatService (service.ChatService)**
   - 处理业务逻辑
   - 管理会话生命周期
   - 处理消息订阅和分发

4. **HTTPServer (server)**
   - 基于 Kratos HTTP 服务器
   - 处理 WebSocket 连接升级
   - 支持并发连接管理

### 工作流程

```
客户端 → WebSocket 连接 → HTTPServer
                              ↓
                          ChatService
                              ↓
                          Session Manager
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
- 使用 Kratos 标准项目结构
- 分层架构：biz（业务逻辑）→ service（服务层）→ server（传输层）
- 完整的注释和文档

### 依赖管理

使用 Go Modules，主要依赖：

- `github.com/go-kratos/kratos/v2` - Kratos 微服务框架
- `github.com/gorilla/websocket` - WebSocket 支持
- `gopkg.in/yaml.v3` - YAML 配置解析

### 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
server:
  http:
    addr: "0.0.0.0:8080"  # 服务器监听地址
    timeout: 5s            # 请求超时时间
```

## 📝 TODO

- [ ] 添加 RESTful API 端点
- [ ] 接入真实的 Kafka 队列
- [ ] 实现真实的 LLM 模型接入
- [ ] 添加连接限流和熔断机制
- [ ] 支持消息持久化
- [ ] 添加 Prometheus 监控指标
- [ ] 实现分布式部署支持
- [ ] 使用 Wire 进行依赖注入
- [ ] 添加 gRPC 支持

## 📄 License

MIT
