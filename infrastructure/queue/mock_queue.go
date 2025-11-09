/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:37
 */
package queue

import (
	"fmt"
	"llm-chat/application/service"
	"llm-chat/domain/valueobject"
	"log"
	"time"
)

// MockQueue 模拟消息队列实现
// 实现 application/service.MessageQueue 接口
type MockQueue struct {
	requests chan *requestMessage
	results  chan *resultMessage
}

type requestMessage struct {
	SessionID valueobject.SessionID
	Content   string
}

type resultMessage struct {
	SessionID valueobject.SessionID
	Chunk     string
	IsLast    bool
}

// NewMockQueue 创建模拟消息队列
func NewMockQueue() service.MessageQueue {
	return &MockQueue{
		requests: make(chan *requestMessage, 1024),
		results:  make(chan *resultMessage, 1024),
	}
}

// StartWorker 启动工作协程
func (mq *MockQueue) StartWorker() error {
	go func() {
		for req := range mq.requests {
			log.Printf("[MockModel] processing: %s -> %s", req.SessionID.Value(), req.Content)
			// 模拟大模型流式输出
			for i := 1; i <= 5; i++ {
				time.Sleep(400 * time.Millisecond)
				mq.results <- &resultMessage{
					SessionID: req.SessionID,
					Chunk:     fmt.Sprintf("chunk %d: %s", i, req.Content),
					IsLast:    i == 5,
				}
			}
		}
	}()
	return nil
}

// PublishRequest 发布请求消息
func (mq *MockQueue) PublishRequest(sessionID valueobject.SessionID, content string) error {
	mq.requests <- &requestMessage{
		SessionID: sessionID,
		Content:   content,
	}
	return nil
}

// SubscribeResults 订阅结果消息
func (mq *MockQueue) SubscribeResults(handler func(sessionID valueobject.SessionID, chunk string, isLast bool)) error {
	go func() {
		for r := range mq.results {
			handler(r.SessionID, r.Chunk, r.IsLast)
		}
	}()
	return nil
}

