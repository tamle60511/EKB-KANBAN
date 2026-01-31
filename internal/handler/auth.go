package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	BaseHandler
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	res, err := h.authService.Login(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Login failed", err)
	}
	return utils.SuccessResponse(c, "Login successful", res)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	res := h.authService.Logout(c.RequestCtx(), req.Token)

	return utils.SuccessResponse(c, "Logout successful", res)
}

func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	userIDRaw := c.Locals("user_id") // Retrieve user_id
	if userIDRaw == nil {
		return utils.InternalErrorResponse(c, "User ID is not set", nil)
	}
	userID, ok := userIDRaw.(int64)
	if !ok {
		return utils.InternalErrorResponse(c, "User ID type assertion failed", nil)
	}
	res, err := h.authService.GetProfile(c.RequestCtx(), userID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Get profile failed", err)
	}
	return utils.SuccessResponse(c, "Get profile successful", res)
}
func (h *AuthHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	auth := router.Group("/auth")
	if len(ms) > 0 {
		for _, m := range ms {
			auth.Use(m)
		}
	}

	auth.Post("/login", h.Login)
	auth.Post("/logout", h.Logout)
	auth.Get("/profile", h.GetProfile)
}
