/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:47
 */
package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func RegisterRoutes(r *gin.Engine, mgr *Manager, mq *MockQueue) {
	r.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.String(400, "upgrade failed: %v", err)
			return
		}

		id := fmt.Sprintf("sess-%d", len(mgr.sessions)+1)
		s := NewSession(id, conn)
		mgr.Register(s)

		go s.WritePump()
		go s.ReadPump(func(sid string, msg []byte) {
			mq.PublishRequest(&Request{SessionID: sid, Content: string(msg)})
		})
	})
}
