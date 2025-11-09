/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package dto

// ChatRequest 聊天请求DTO
type ChatRequest struct {
	SessionID string `json:"session_id"`
	Content   string `json:"content"`
}

// CreateSessionRequest 创建会话请求DTO
type CreateSessionRequest struct {
	SessionID string `json:"session_id,omitempty"` // 可选，如果不提供则自动生成
}

