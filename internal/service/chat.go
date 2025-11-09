package service

import (
	"fmt"
	"llm-chat/internal/biz"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// ChatService 聊天服务
type ChatService struct {
	manager *biz.Manager
	queue   biz.Queue
	counter int64
}

// NewChatService 创建聊天服务
func NewChatService(manager *biz.Manager, queue biz.Queue) *ChatService {
	cs := &ChatService{
		manager: manager,
		queue:   queue,
	}
	
	// 启动模型结果订阅
	cs.startResultSubscription()
	
	return cs
}

// startResultSubscription 启动结果订阅
func (cs *ChatService) startResultSubscription() {
	cs.queue.SubscribeResults(func(res *biz.Result) {
		if sess, ok := cs.manager.Get(res.SessionID); ok {
			select {
			case sess.Send <- []byte(res.Chunk):
			default:
				log.Printf("[WARN] session %s send buffer full, drop chunk", res.SessionID)
			}
		}
	})
}

// HandleWebSocket 处理 WebSocket 连接
func (cs *ChatService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源，生产环境应该检查
		},
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("upgrade failed: %v", err), http.StatusBadRequest)
		return
	}
	
	// 生成会话ID
	id := fmt.Sprintf("sess-%d", atomic.AddInt64(&cs.counter, 1))
	session := biz.NewSession(id, conn)
	
	// 注册会话
	cs.manager.Register(session)
	
	// 启动读写协程
	go session.WritePump()
	go session.ReadPump(func(sid string, msg []byte) {
		cs.queue.PublishRequest(&biz.Request{
			SessionID: sid,
			Content:   string(msg),
		})
	})
	
	// 会话结束时注销
	go func() {
		<-session.CloseCh
		cs.manager.Unregister(id)
		log.Printf("[ChatService] session unregistered: %s", id)
	}()
	
	log.Printf("[ChatService] new session registered: %s", id)
}

