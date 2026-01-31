package handler

import (
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type CopmaHandler struct {
	BaseHandler
	CopmaService service.CopmaService
}

func NewCopmaHandler(copmaService service.CopmaService) *CopmaHandler {
	return &CopmaHandler{
		CopmaService: copmaService,
	}
}

func (h *CopmaHandler) GetCopma(c fiber.Ctx) error {
	res, err := h.CopmaService.GetCopma(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get COPMA data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved COPMA data", res)
}

func (h *CopmaHandler) GetChannel(c fiber.Ctx) error {
	res, err := h.CopmaService.GetChannel(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Channel data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Channel data", res)
}
func (h *CopmaHandler) GetTypes(c fiber.Ctx) error {
	res, err := h.CopmaService.GetTypes(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Types data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Types data", res)
}
func (h *CopmaHandler) GetRegion(c fiber.Ctx) error {
	res, err := h.CopmaService.GetRegion(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Region data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Region data", res)
}
func (h *CopmaHandler) GetCountry(c fiber.Ctx) error {
	res, err := h.CopmaService.GetCountry(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Country data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Country data", res)
}
func (h *CopmaHandler) GetRoute(c fiber.Ctx) error {
	res, err := h.CopmaService.GetRoute(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Route data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Route data", res)
}
func (h *CopmaHandler) GetSaleDept(c fiber.Ctx) error {
	res, err := h.CopmaService.GetSaleDept(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Sale Department data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Sale Department data", res)
}

func (h *CopmaHandler) GetSaleWorkshop(c fiber.Ctx) error {
	res, err := h.CopmaService.GetSaleWorkshop(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Sale Workshop data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Sale Workshop data", res)
}
func (h *CopmaHandler) GetSaleItem(c fiber.Ctx) error {
	res, err := h.CopmaService.GetSaleItem(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Sale Item data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Sale Item data", res)
}

func (h *CopmaHandler) GetSaleWarehouse(c fiber.Ctx) error {
	res, err := h.CopmaService.GetSaleWarehouse(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Sale Warehouse data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Sale Warehouse data", res)
}

func (h *CopmaHandler) GetSaleMoney(c fiber.Ctx) error {
	res, err := h.CopmaService.GetSaleMoney(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get Sale Money data", err)
	}
	return utils.SuccessResponse(c, "Successfully retrieved Sale Money data", res)
}
func (h *CopmaHandler) SearchCopma(c fiber.Ctx) error {
	search := c.Query("search", "")
	res, err := h.CopmaService.SearchCopma(c.RequestCtx(), search)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to search COPMA data", err)
	}
	return utils.SuccessResponse(c, "Successfully searched COPMA data", res)
}
func (h *CopmaHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	r := router.Group("/copma")
	for _, m := range ms {
		r.Use(m)
	}
	r.Get("/", h.GetCopma)
	r.Get("/channel", h.GetChannel)
	r.Get("/types", h.GetTypes)
	r.Get("/region", h.GetRegion)
	r.Get("/country", h.GetCountry)
	r.Get("/route", h.GetRoute)
	r.Get("/saledept", h.GetSaleDept)
	r.Get("/saleworkshop", h.GetSaleWorkshop)
	r.Get("/saleitem", h.GetSaleItem)
	r.Get("/salewarehouse", h.GetSaleWarehouse)
	r.Get("/salemoney", h.GetSaleMoney)
	r.Get("/search", h.SearchCopma)
}
