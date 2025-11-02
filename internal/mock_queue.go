/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:37
 */
package internal

import (
	"fmt"
	"log"
	"time"
)

// 模拟一个Kafka/Redis流队列
type MockQueue struct {
	requests chan *Request
	results  chan *Result
}

type Request struct {
	SessionID string
	Content   string
}

type Result struct {
	SessionID string
	Chunk     string
	IsLast    bool
}

func NewMockQueue() *MockQueue {
	return &MockQueue{
		requests: make(chan *Request, 1024),
		results:  make(chan *Result, 1024),
	}
}

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

// 模拟异步发布请求
func (mq *MockQueue) PublishRequest(req *Request) {
	mq.requests <- req
}

func (mq *MockQueue) SubscribeResults(deliver func(*Result)) {
	go func() {
		for r := range mq.results {
			deliver(r)
		}
	}()
}
