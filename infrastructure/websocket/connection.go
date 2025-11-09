/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package websocket

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

// Connection WebSocket连接封装
// 处理WebSocket的技术细节
type Connection struct {
	Conn     *websocket.Conn
	CloseCh  chan struct{}
	MessageCh chan []byte
}

// NewConnection 创建WebSocket连接
func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		Conn:      conn,
		CloseCh:   make(chan struct{}),
		MessageCh: make(chan []byte, 128),
	}
}

// ReadPump WebSocket读协程
// 只用于接收用户消息和维持连接
func (c *Connection) ReadPump(handleMessage func([]byte)) {
	defer func() {
		log.Printf("[WebSocket] closed readPump")
		close(c.CloseCh)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(1 << 20)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("[WebSocket] read err: %v", err)
			break
		}
		handleMessage(msg)
	}
}

// HeartbeatPump WebSocket心跳协程
// 只用于维持连接和心跳
func (c *Connection) HeartbeatPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[WebSocket] heartbeat err: %v", err)
				return
			}
		case <-c.CloseCh:
			return
		}
	}
}

// Close 关闭连接
func (c *Connection) Close() error {
	close(c.CloseCh)
	return c.Conn.Close()
}

