# llm_chat
一个最小可运行（standalone）高并发LLM聊天 demo

## 核心模块

* **Gin + gorilla/websocket** 实现 RESTful & WebSocket 接入
* 每个会话（Session）拥有独立 goroutine，使用 channel 管理 I/O
* 模拟 **Kafka/Redis 队列**（用 Go channel stub 实现消息解耦）
* 流式返回（定时推送生成内容）
* 支持 10 万连接的基础结构（单机 demo 可压测几千连接）

---

## 🧱 项目结构

```
chat-demo/
├── go.mod
├── main.go
├── internal/
│   ├── session.go        # Session、Manager
│   ├── ws_handler.go     # Gin + WebSocket handler
│   └── mock_queue.go     # 模拟Kafka/Redis的消息队列
```

---

## ✅ 完整代码

> 直接复制保存成一个目录后运行 `go run .` 即可。
> 运行后访问：[http://localhost:8080/ws](http://localhost:8080/ws)
> （建议用 WebSocket 客户端连接测试，如 Postman、或浏览器 console）
