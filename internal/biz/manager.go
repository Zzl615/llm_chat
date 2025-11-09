package biz

import (
	"sync"
)

// Manager 维护所有活跃 session
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewManager 创建新的管理器
func NewManager() *Manager {
	return &Manager{sessions: make(map[string]*Session)}
}

// Register 注册会话
func (m *Manager) Register(s *Session) {
	m.mu.Lock()
	m.sessions[s.ID] = s
	m.mu.Unlock()
}

// Unregister 注销会话
func (m *Manager) Unregister(id string) {
	m.mu.Lock()
	delete(m.sessions, id)
	m.mu.Unlock()
}

// Get 获取会话
func (m *Manager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	s, ok := m.sessions[id]
	m.mu.RUnlock()
	return s, ok
}

// Count 返回当前活跃会话数
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}
