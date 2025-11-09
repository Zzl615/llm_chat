/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package repository

import (
	"llm-chat/domain/entity"
	"llm-chat/domain/repository"
	"llm-chat/domain/valueobject"
	"sync"
)

// SessionRepositoryImpl 会话仓储实现
// 实现 domain/repository.SessionRepository 接口
type SessionRepositoryImpl struct {
	mu       sync.RWMutex
	sessions map[string]*entity.Session
}

// NewSessionRepositoryImpl 创建会话仓储实现
func NewSessionRepositoryImpl() repository.SessionRepository {
	return &SessionRepositoryImpl{
		sessions: make(map[string]*entity.Session),
	}
}

// Save 保存会话
func (r *SessionRepositoryImpl) Save(session *entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID.Value()] = session
	return nil
}

// FindByID 根据ID查找会话
func (r *SessionRepositoryImpl) FindByID(id valueobject.SessionID) (*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[id.Value()]
	if !ok {
		return nil, nil
	}
	return session, nil
}

// Delete 删除会话
func (r *SessionRepositoryImpl) Delete(id valueobject.SessionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id.Value())
	return nil
}

// FindAll 查找所有会话
func (r *SessionRepositoryImpl) FindAll() ([]*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sessions := make([]*entity.Session, 0, len(r.sessions))
	for _, session := range r.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// Count 统计会话数量
func (r *SessionRepositoryImpl) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions)
}

