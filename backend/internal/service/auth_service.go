package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/taskmanager/backend/internal/config"
	"github.com/taskmanager/backend/internal/models"
	"github.com/taskmanager/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type AuthService interface {
	Register(email, password, name string) (*models.User, error)
	Login(email, password string) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	Logout(userID uuid.UUID) error
	ValidateAccessToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	config    *config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, cfg *config.JWTConfig) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		config:    cfg,
	}
}

func (s *authService) Register(email, password, name string) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(user)
}

func (s *authService) RefreshToken(refreshToken string) (*TokenPair, error) {
	storedToken, err := s.tokenRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if storedToken.IsExpired() {
		s.tokenRepo.Delete(storedToken.ID)
		return nil, ErrTokenExpired
	}

	s.tokenRepo.Delete(storedToken.ID)

	return s.generateTokenPair(&storedToken.User)
}

func (s *authService) Logout(userID uuid.UUID) error {
	return s.tokenRepo.DeleteByUserID(userID)
}

func (s *authService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *authService) generateTokenPair(user *models.User) (*TokenPair, error) {
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessExpiry.Seconds()),
	}, nil
}

func (s *authService) generateAccessToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

func (s *authService) generateRefreshToken(user *models.User) (string, error) {
	tokenID := uuid.New()
	expiresAt := time.Now().Add(s.config.RefreshExpiry)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
		ID:        tokenID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}

	refreshToken := &models.RefreshToken{
		Token:     signedToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return "", err
	}

	return signedToken, nil
}
