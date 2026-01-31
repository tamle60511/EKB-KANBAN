package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type DepartmentHandler struct {
	BaseHandler
	departmentService service.DepartmentService
}

func NewDepartmentHandler(departmentService service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

func (d *DepartmentHandler) Create(c fiber.Ctx) error {
	var req dto.DepartmenCreateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := d.departmentService.Create(c.RequestCtx(), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to create ddepartment", err)
	}
	return utils.SuccessResponse(c, "create department success", nil)
}

func (d *DepartmentHandler) GetByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid department id", err)
	}
	dept, err := d.departmentService.GetByID(c.RequestCtx(), uint(id))
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get department by id", err)
	}
	if dept == nil {
		return utils.SuccessResponse(c, "department not found", nil)
	}
	return utils.SuccessResponse(c, "get department success", dept)
}

func (d *DepartmentHandler) GetAll(c fiber.Ctx) error {
	depts, err := d.departmentService.GetAll(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get all departments", err)
	}
	return utils.SuccessResponse(c, "get all departments success", depts)
}

func (d *DepartmentHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid department id", err)
	}
	var req dto.DepartmentUpdateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := d.departmentService.Update(c.RequestCtx(), uint(id), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update department", err)
	}
	return utils.SuccessResponse(c, "update department success", nil)
}
func (d *DepartmentHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid department id", err)
	}
	if err := d.departmentService.Delete(c.RequestCtx(), uint(id)); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete department", err)
	}
	return utils.SuccessResponse(c, "delete department success", nil)
}
func (d *DepartmentHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	department := router.Group("/departments")
	if len(ms) > 0 {
		for _, m := range ms {
			department.Use(m)
		}
	}
	department.Post("/", d.Create)
	department.Get("/", d.GetAll)
	department.Get("/:id", d.GetByID)
	department.Put("/:id", d.Update)
	department.Delete("/:id", d.Delete)
}
