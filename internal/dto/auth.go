package dto

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	DepartmentID int64  `json:"department_id"`
	Role         string `json:"role"`
}

type ValidateTokenRequest struct {
	Token string `json:"token" `
}

type ValidateTokenResponse struct {
	Valid bool `json:"valid"`
}

type LogoutRequest struct {
	Token string `json:"token" `
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type TokenClaims struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	DepartmentID int64  `json:"department_id"`
	Exp          int64  `json:"exp,omitempty"`
	Role         string `json:"role,omitempty"`
	jwt.RegisteredClaims
}
