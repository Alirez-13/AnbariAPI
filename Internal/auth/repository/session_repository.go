package repository

import (
	"AnbariAPI/Internal/auth/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *domain.Session) error
	FindByID(id uuid.UUID) (*domain.Session, error)
	FindByUserID(userID uint) ([]domain.Session, error)
	Deactivate(id uuid.UUID) error
	DeactivateAllByUserID(userID uint) error
	DeleteExpired() error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) Create(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) FindByID(id uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	result := r.db.Where("id = ? AND is_active = ?", id, true).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *sessionRepository) FindByUserID(userID uint) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions).Error
	return sessions, err
}

func (r *sessionRepository) Deactivate(id uuid.UUID) error {
	return r.db.Model(&domain.Session{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *sessionRepository) DeactivateAllByUserID(userID uint) error {
	return r.db.Model(&domain.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

func (r *sessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&domain.Session{}).Error
}
