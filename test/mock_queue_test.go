/**
 * @Author: Noaghzil
 * @Date:   2025-11-02
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02
 */
package test

import (
	"fmt"
	"llm-chat/internal"
	"sync"
	"testing"
	"time"
)

func TestNewMockQueue(t *testing.T) {
	mq := internal.NewMockQueue()
	if mq == nil {
		t.Fatal("NewMockQueue returned nil")
	}
}

func TestPublishRequest(t *testing.T) {
	mq := internal.NewMockQueue()

	req := &internal.Request{
		SessionID: "test-session-1",
		Content:   "Hello, world!",
	}

	// 测试发布请求不会阻塞
	done := make(chan bool)
	go func() {
		mq.PublishRequest(req)
		done <- true
	}()

	select {
	case <-done:
		// 成功发布
	case <-time.After(1 * time.Second):
		t.Fatal("PublishRequest blocked or timed out")
	}
}

func TestStartMockModelWorker(t *testing.T) {
	mq := internal.NewMockQueue()
	mq.StartMockModelWorker()

	// 等待 worker 启动
	time.Sleep(100 * time.Millisecond)

	// 发布一个请求
	mq.PublishRequest(&internal.Request{
		SessionID: "test-session-2",
		Content:   "Test message",
	})

	// 订阅结果
	results := make([]*internal.Result, 0)
	done := make(chan bool)

	mq.SubscribeResults(func(res *internal.Result) {
		results = append(results, res)
		if res.IsLast {
			done <- true
		}
	})

	// 等待所有结果（5个chunk）
	select {
	case <-done:
		// 检查是否收到5个结果
		if len(results) != 5 {
			t.Fatalf("Expected 5 results, got %d", len(results))
		}
		// 检查最后一个结果的 IsLast 标志
		if !results[4].IsLast {
			t.Error("Last result should have IsLast = true")
		}
		// 检查 SessionID 是否正确
		if results[0].SessionID != "test-session-2" {
			t.Errorf("Expected SessionID 'test-session-2', got '%s'", results[0].SessionID)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for results")
	}
}

func TestMultipleRequests(t *testing.T) {
	mq := internal.NewMockQueue()
	mq.StartMockModelWorker()

	// 等待 worker 启动
	time.Sleep(100 * time.Millisecond)

	// 收集所有结果（在发布请求之前先订阅）
	resultMap := make(map[string][]*internal.Result)
	var mu sync.Mutex
	completedSessions := make(map[string]bool)
	done := make(chan string, 3) // 3个session，每个5个chunk，最后一个是IsLast

	mq.SubscribeResults(func(res *internal.Result) {
		mu.Lock()
		resultMap[res.SessionID] = append(resultMap[res.SessionID], res)
		if res.IsLast && !completedSessions[res.SessionID] {
			completedSessions[res.SessionID] = true
			done <- res.SessionID
		}
		mu.Unlock()
	})

	// 发布多个请求
	sessionIDs := []string{"session-1", "session-2", "session-3"}
	for _, sid := range sessionIDs {
		mq.PublishRequest(&internal.Request{
			SessionID: sid,
			Content:   "Message for " + sid,
		})
	}

	// 等待所有session完成（每个请求需要5 * 400ms = 2秒，加上一些缓冲）
	timeout := time.After(10 * time.Second)
	completed := make(map[string]bool)

	for len(completed) < len(sessionIDs) {
		select {
		case sid := <-done:
			completed[sid] = true
		case <-timeout:
			t.Fatalf("Timeout waiting for all sessions to complete. Completed: %v", completed)
		}
	}

	// 检查每个 session 都应该收到 5 个结果
	for _, sid := range sessionIDs {
		mu.Lock()
		results, ok := resultMap[sid]
		mu.Unlock()
		if !ok {
			t.Errorf("No results received for session %s", sid)
			continue
		}
		if len(results) != 5 {
			t.Errorf("Session %s expected 5 results, got %d", sid, len(results))
		}
	}
}

func TestConcurrentPublish(t *testing.T) {
	mq := internal.NewMockQueue()
	mq.StartMockModelWorker()

	// 等待 worker 启动
	time.Sleep(100 * time.Millisecond)

	// 并发发布多个请求
	const numRequests = 10
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			mq.PublishRequest(&internal.Request{
				SessionID: fmt.Sprintf("concurrent-session-%d", id),
				Content:   fmt.Sprintf("Message %d", id),
			})
			done <- true
		}(i)
	}

	// 等待所有发布完成
	for i := 0; i < numRequests; i++ {
		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatalf("Timeout waiting for request %d to be published", i)
		}
	}
}
