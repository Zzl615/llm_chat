/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package entity

import (
	"llm-chat/domain/valueobject"
	"time"
)

// Session 会话领域实体
// 只包含业务属性，不包含技术实现细节（如WebSocket连接）
type Session struct {
	ID        valueobject.SessionID
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

// NewSession 创建新的会话实体
func NewSession(id valueobject.SessionID) *Session {
	now := time.Now()
	return &Session{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}
}

// Activate 激活会话
func (s *Session) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate 停用会话
func (s *Session) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// IsValid 检查会话是否有效
func (s *Session) IsValid() bool {
	return s.IsActive
}

