/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:32
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:17:51
 */
package main

import (
	"llm-chat/internal"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	manager := internal.NewManager()
	queue := internal.NewMockQueue()

	// 启动模拟的模型推理流
	queue.StartMockModelWorker()

	// 模型结果订阅：将结果投递给对应 session
	queue.SubscribeResults(func(res *internal.Result) {
		if sess, ok := manager.Get(res.SessionID); ok {
			select {
			case sess.Send <- []byte(res.Chunk):
			default:
				log.Printf("[WARN] session %s send buffer full, drop chunk", res.SessionID)
			}
		}
	})

	internal.RegisterRoutes(r, manager, queue)

	log.Println("🚀 Chat demo server started at :8080")
	r.Run(":8080")
}
