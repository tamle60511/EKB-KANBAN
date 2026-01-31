// handler/sale_copi04_handler.go

package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type SaleCopi04Handler struct {
	BaseHandler
	SaleCopi04Service service.SaleCopi04Service
}

func NewSaleCopi04Handler(saleCopi04Service service.SaleCopi04Service) *SaleCopi04Handler {
	return &SaleCopi04Handler{
		SaleCopi04Service: saleCopi04Service,
	}
}

// ✅ Get single record
func (h *SaleCopi04Handler) GetCopi04(c fiber.Ctx) error {
	var req dto.SaleCopi04Req
	if err := c.Bind().URI(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid URI params", err)
	}

	res, err := h.SaleCopi04Service.GetCopi04(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get COPI04 data", err)
	}

	return utils.SuccessResponse(c, "Successfully retrieved COPI04 data", res)
}

// ✅ Get all records
func (h *SaleCopi04Handler) GetAllCopi04(c fiber.Ctx) error {
	res, err := h.SaleCopi04Service.GetAllCopi04(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get all COPI04 data", err)
	}

	return utils.SuccessResponse(c, "Successfully retrieved all COPI04 data", res)
}

func (h *SaleCopi04Handler) Create(c fiber.Ctx) error {
	var input dto.SaleCopi04Create
	if err := c.Bind().Body(&input); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}

	// Validate
	if input.ME001 == "" {
		return utils.BadRequestResponse(c, "ME001 is required", nil)
	}


	creator := c.Locals("username").(string)
	company := "CQS_VN_2025"

	// Create
	if err := h.SaleCopi04Service.Create(c.RequestCtx(), input, creator, company); err != nil {
		// Handle duplicate key error
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "PRIMARY KEY") {
			return utils.BadRequestResponse(c, "ME001 already exists", nil)
		}
		return utils.InternalErrorResponse(c, "Failed to create COPI04 data", err)
	}

	// ✅ Fetch and return created data
	req := dto.SaleCopi04Req{ME001: input.ME001}
	result, err := h.SaleCopi04Service.GetCopi04(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Created but failed to fetch data", err)
	}

	return utils.SuccessResponse(c, "Successfully created COPI04 data", result)
}

func (h *SaleCopi04Handler) Update(c fiber.Ctx) error {
	var req dto.SaleCopi04Req
	if err := c.Bind().URI(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid URI params", err)
	}

	var input dto.SaleCopi04Update
	if err := c.Bind().Body(&input); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}

	modifier := c.Locals("username").(string)
	company := "CQS_VN_2025"

	if err := h.SaleCopi04Service.Update(c.RequestCtx(), req.ME001, input, modifier, company); err != nil {
		return utils.InternalErrorResponse(c, "Failed to update COPI04 data", err)
	}

	result, err := h.SaleCopi04Service.GetCopi04(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Updated but failed to fetch data", err)
	}

	return utils.SuccessResponse(c, "Successfully updated COPI04 data", result)
}

// ✅ Delete
func (h *SaleCopi04Handler) Delete(c fiber.Ctx) error {
	var req dto.SaleCopi04Req
	if err := c.Bind().URI(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid URI params", err)
	}

	if err := h.SaleCopi04Service.Delete(c.RequestCtx(), req.ME001); err != nil {
		return utils.InternalErrorResponse(c, "Failed to delete COPI04 data", err)
	}

	return utils.SuccessResponse(c, "Successfully deleted COPI04 data", nil)
}

// ✅ Setup routes
func (h *SaleCopi04Handler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	saleCopi04 := router.Group("/sale-copi04")
	saleCopi04.Get("/", h.GetAllCopi04) // List all
	saleCopi04.Get("/:id", h.GetCopi04) // Get by ID
	saleCopi04.Post("/", h.Create)      // Create
	saleCopi04.Put("/:id", h.Update)    // Update
	saleCopi04.Delete("/:id", h.Delete) // Delete
}
