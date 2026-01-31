package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type MenuHandler struct {
	BaseHandler
	menuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

func (m *MenuHandler) CreateMenu(c fiber.Ctx) error {
	var menus dto.MenuCreateReq
	if err := c.Bind().Body(&menus); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := m.menuService.Create(c.RequestCtx(), menus); err != nil {
		return utils.InternalErrorResponse(c, "failed to create menus", err)
	}
	return utils.SuccessResponse(c, "create menus success", nil)
}
func (m *MenuHandler) UpdateMenu(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID params")
	}
	var req dto.MenuUpdateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := m.menuService.Update(c.RequestCtx(), int64(id), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update report", err)
	}
	return utils.SuccessResponse(c, "Update to report success", nil)
}
func (m *MenuHandler) DeleteMenu(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID params")
	}
	if err := m.menuService.Delete(c.RequestCtx(), id); err != nil {
		return utils.BadRequestResponse(c, "failed to delete", err)
	}
	return utils.SuccessResponse(c, "delete menus success", nil)
}
func (m *MenuHandler) GetMenu(c fiber.Ctx) error {
	req := dto.MenuDetailReq{}
	if err := c.Bind().Query(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid query body", err)
	}
	data, err := m.menuService.GetMenu(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get menu", err)
	}

	return utils.SuccessResponse(c, "get menu success", data)
}
func (m *MenuHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID params")
	}
	data, err := m.menuService.GetByID(c.RequestCtx(), id)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get menu %w", err)
	}
	return utils.SuccessResponse(c, "get menu success", data)
}
func (m *MenuHandler) GetAllMenu(c fiber.Ctx) error {
	req := dto.MenuDetailReq{}
	if err := c.Bind().Query(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid query body", err)
	}
	data, err := m.menuService.GetAll(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get menu %w", err)
	}
	return utils.SuccessResponse(c, "get menu success", data)
}
func (m *MenuHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	menus := router.Group("/menus")
	if len(ms) > 0 {
		for _, m := range ms {
			menus.Use(m)
		}
	}
	menus.Post("/", m.CreateMenu)
	menus.Get("/:id", m.GetMenu)
	menus.Get("/by-id/:id", m.GetByID)
	menus.Put("/:id", m.UpdateMenu)
	menus.Delete("/:id", m.DeleteMenu)
	menus.Get("/", m.GetAllMenu)
}
