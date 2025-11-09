/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package dto

// ChatResponse 聊天响应DTO
type ChatResponse struct {
	SessionID string `json:"session_id"`
	Chunk     string `json:"chunk"`
	IsLast    bool   `json:"is_last"`
}

// SessionResponse 会话响应DTO
type SessionResponse struct {
	SessionID string `json:"session_id"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// ErrorResponse 错误响应DTO
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

