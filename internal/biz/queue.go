package biz

import (
	"fmt"
	"log"
	"time"
)

// Queue 消息队列接口
type Queue interface {
	PublishRequest(req *Request)
	SubscribeResults(deliver func(*Result))
	StartMockModelWorker()
}

// MockQueue 模拟一个Kafka/Redis流队列
type MockQueue struct {
	requests chan *Request
	results  chan *Result
}

// Request 请求消息
type Request struct {
	SessionID string
	Content   string
}

// Result 结果消息
type Result struct {
	SessionID string
	Chunk     string
	IsLast    bool
}

// NewMockQueue 创建新的模拟队列
func NewMockQueue() *MockQueue {
	return &MockQueue{
		requests: make(chan *Request, 1024),
		results:  make(chan *Result, 1024),
	}
}

// StartMockModelWorker 启动模拟的模型推理工作协程
func (mq *MockQueue) StartMockModelWorker() {
	go func() {
		for req := range mq.requests {
			log.Printf("[MockModel] processing: %s -> %s", req.SessionID, req.Content)
			// 模拟大模型流式输出
			for i := 1; i <= 5; i++ {
				time.Sleep(400 * time.Millisecond)
				mq.results <- &Result{
					SessionID: req.SessionID,
					Chunk:     fmt.Sprintf("chunk %d: %s", i, req.Content),
					IsLast:    i == 5,
				}
			}
		}
	}()
}

// PublishRequest 发布请求到队列
func (mq *MockQueue) PublishRequest(req *Request) {
	mq.requests <- req
}

// SubscribeResults 订阅结果
func (mq *MockQueue) SubscribeResults(deliver func(*Result)) {
	go func() {
		for r := range mq.results {
			deliver(r)
		}
	}()
}

