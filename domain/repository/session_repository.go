/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package repository

import (
	"llm-chat/domain/entity"
	"llm-chat/domain/valueobject"
)

// SessionRepository 会话仓储接口
// 定义在领域层，由基础设施层实现（依赖倒置原则）
type SessionRepository interface {
	// Save 保存会话
	Save(session *entity.Session) error

	// FindByID 根据ID查找会话
	FindByID(id valueobject.SessionID) (*entity.Session, error)

	// Delete 删除会话
	Delete(id valueobject.SessionID) error

	// FindAll 查找所有会话
	FindAll() ([]*entity.Session, error)

	// Count 统计会话数量
	Count() int
}

