// # SINGLE REASON: Persist auth sessions with GORM.
package infrastructure

import (
	"errors"
	"time"

	"AnbariAPI/internal/auth/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormSessionRepository struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

func (r *GormSessionRepository) Create(session *domain.Session) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	return r.db.Create(session).Error
}

func (r *GormSessionRepository) FindByID(id uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	result := r.db.Where("id = ? AND is_active = ?", id, true).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return &session, nil
}

func (r *GormSessionRepository) FindByUserID(userID uint) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions).Error
	return sessions, err
}

func (r *GormSessionRepository) Deactivate(id uuid.UUID) error {
	return r.db.Model(&domain.Session{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *GormSessionRepository) DeactivateAllByUserID(userID uint) error {
	return r.db.Model(&domain.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

func (r *GormSessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&domain.Session{}).Error
}
