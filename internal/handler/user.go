package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	BaseHandler
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) Create(c fiber.Ctx) error {
	var req dto.UserCreateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := u.userService.Create(c.RequestCtx(), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to create users", err)
	}
	return utils.SuccessResponse(c, "create user success", nil)
}

func (u *UserHandler) GetAll(c fiber.Ctx) error {
	users, err := u.userService.GetAll(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get users", err)
	}
	return utils.SuccessResponse(c, "get users success", users)
}

func (u *UserHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID", err)
	}
	var req dto.UserUpdateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := u.userService.Update(c.RequestCtx(), id, req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update user", err)
	}
	return utils.SuccessResponse(c, "update user success", nil)
}
func (u *UserHandler) UpdatePassword(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID", err)
	}
	var req dto.UpdatePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if req.NewPassword == "" {
		return utils.BadRequestResponse(c, "New password cannot be empty", nil)
	}
	if err := u.userService.UpdatePassword(c.RequestCtx(), int64(id), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update user password", err)
	}
	return utils.SuccessResponse(c, "update user password success", nil)
}
func (u *UserHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID", err)
	}

	if err := u.userService.Delete(c.RequestCtx(), id); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete user", err)
	}

	return utils.SuccessResponse(c, "delete user success", nil)
}
func (u *UserHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID", err)
	}
	req := dto.UserDetailReq{UserID: int64(id)}
	user, err := u.userService.GetByID(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get user", err)
	}
	return utils.SuccessResponse(c, "get user success", user)
}
func (u *UserHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	users := router.Group("/users")
	if len(ms) > 0 {
		for _, m := range ms {
			users.Use(m)
		}
	}
	users.Post("/", u.Create)
	users.Get("/", u.GetAll)
	users.Put("/:id", u.Update)
	users.Delete("/:id", u.Delete)
	users.Put("/:id/password", u.UpdatePassword)
	users.Get("/:id", u.GetByID)
}
