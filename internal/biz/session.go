package biz

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 5 * time.Second
	pongWait   = 30 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// Session 表示一个用户会话
type Session struct {
	ID      string
	Conn    *websocket.Conn
	Send    chan []byte // 发送缓冲区
	CloseCh chan struct{}
}

// NewSession 创建新的会话
func NewSession(id string, conn *websocket.Conn) *Session {
	return &Session{
		ID:      id,
		Conn:    conn,
		Send:    make(chan []byte, 128),
		CloseCh: make(chan struct{}),
	}
}

// ReadPump WebSocket读协程
func (s *Session) ReadPump(handleRequest func(string, []byte)) {
	defer func() {
		log.Printf("[session %s] closed readPump", s.ID)
		close(s.CloseCh)
		s.Conn.Close()
	}()

	s.Conn.SetReadLimit(1 << 20)
	s.Conn.SetReadDeadline(time.Now().Add(pongWait))
	s.Conn.SetPongHandler(func(string) error {
		s.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := s.Conn.ReadMessage()
		if err != nil {
			log.Printf("[session %s] read err: %v", s.ID, err)
			break
		}
		handleRequest(s.ID, msg)
	}
}

// WritePump WebSocket写协程
func (s *Session) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-s.Send:
			s.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				s.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := s.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[session %s] write err: %v", s.ID, err)
				return
			}
		case <-ticker.C:
			s.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-s.CloseCh:
			return
		}
	}
}
