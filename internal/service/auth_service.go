package service

import (
	"context"
	"cqs-kanban/config"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
	"cqs-kanban/internal/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*dto.TokenClaims, error)
	GenerateToken(user *models.User) (string, error)
	Logout(ctx context.Context, token string) error
	GetProfile(ctx context.Context, userID int64) (*dto.UserDetailRes, error)
}

type authService struct {
	userRepo repository.UserRepo
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepo, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		log.Printf("Password hash mismatch for user: %s", req.Username)
		return nil, errors.New("invalid username or password")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:        token,
		UserID:       user.ID,
		Username:     user.Username,
		DepartmentID: user.DepartmentID,
		Role:         user.Role,
	}, nil
}
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*dto.TokenClaims, error) {
	claims := &dto.TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
func (s *authService) GenerateToken(user *models.User) (string, error) {

	expirationTime := time.Now().Add(s.config.GetJWTExpiry())
	claims := dto.TokenClaims{
		UserID:       user.ID,
		Username:     user.Username,
		DepartmentID: user.DepartmentID,
		Role:         user.Role,
		Exp:          expirationTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	_, err := s.ValidateToken(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}

	// TODO: Implement token blacklisting or other logout logic
	// Ví dụ: Lưu token vào danh sách bị chặn
	// return s.tokenBlacklistRepo.Add(ctx, token)

	return nil
}

func (s *authService) GetProfile(ctx context.Context, userID int64) (*dto.UserDetailRes, error) {
	return s.userRepo.GetUser(ctx, dto.UserDetailReq{UserID: userID})
}
