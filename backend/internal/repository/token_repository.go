package repository

import (
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	Delete(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Preload("User").First(&refreshToken, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *tokenRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.RefreshToken{}, "id = ?", id).Error
}

func (r *tokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&models.RefreshToken{}, "user_id = ?", userID).Error
}

func (r *tokenRepository) DeleteExpired() error {
	return r.db.Delete(&models.RefreshToken{}, "expires_at < NOW()").Error
}
