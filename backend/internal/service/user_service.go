package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id uuid.UUID, name, avatar string) (*models.User, error)
	UpdatePassword(id uuid.UUID, oldPassword, newPassword string) error
	Delete(id uuid.UUID) error
	List(page, limit int) ([]models.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) Update(id uuid.UUID, name, avatar string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if name != "" {
		user.Name = name
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdatePassword(id uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

func (s *userService) List(page, limit int) ([]models.User, int64, error) {
	return s.userRepo.List(page, limit)
}
