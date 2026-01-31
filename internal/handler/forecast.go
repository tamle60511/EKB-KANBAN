package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ForecastHandler struct {
	BaseHandler
	forecastService service.ForecastService
}

func NewForecastHandler(forecastService service.ForecastService) *ForecastHandler {
	return &ForecastHandler{
		forecastService: forecastService,
	}
}


func (h *ForecastHandler) GetForecast(c fiber.Ctx) error {
	
	fromDateStr := c.Query("FromDate")
	toDateStr := c.Query("ToDate")

	if fromDateStr == "" || toDateStr == "" {
		return utils.BadRequestResponse(c, "FromDate and ToDate are required (format: YYYYMMDD)", nil)
	}

	fromDate, err := parseCompactDate(fromDateStr)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid FromDate format, expected YYYYMMDD", err.Error())
	}

	toDate, err := parseCompactDate(toDateStr)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid ToDate format, expected YYYYMMDD", err.Error())
	}

	req := dto.ReportReq{
		FromDate: &fromDate,
		ToDate:   &toDate,
	}
	if err := c.Bind().URI(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid URI params", err.Error())
	}
	reqCtx, err := h.extractRequestContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}

	forecasts, err := h.forecastService.GetForecast(c.RequestCtx(), reqCtx, &req, c)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get forecast", err)
	}

	return utils.SuccessResponse(c, "get forecast success", forecasts)
}

func (h *ForecastHandler) ExportReport(c fiber.Ctx) error {
	fromDateStr := c.Query("FromDate")
	toDateStr := c.Query("ToDate")

	if fromDateStr == "" || toDateStr == "" {
		return utils.BadRequestResponse(c, "FromDate and ToDate are required (format: YYYYMMDD)", nil)
	}

	fromDate, err := parseCompactDate(fromDateStr)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid FromDate format, expected YYYYMMDD", err.Error())
	}

	toDate, err := parseCompactDate(toDateStr)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid ToDate format, expected YYYYMMDD", err.Error())
	}

	req := dto.ReportReq{
		FromDate: &fromDate,
		ToDate:   &toDate,
	}

	reqCtx, err := h.extractRequestContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}

	reportData, err := h.forecastService.ExportReport(c.RequestCtx(), reqCtx, &req, c)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to export report", err)
	}

	if fileBytes, ok := reportData.FileDetal.([]byte); ok {
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, reportData.FileName))
		c.Set("Content-Length", strconv.Itoa(len(fileBytes)))
		return c.Send(fileBytes)
	}

	return utils.InternalErrorResponse(c, "Invalid file data", "reportData.FileDetal is not []byte")
}

func parseCompactDate(dateStr string) (time.Time, error) {
	if len(dateStr) != 8 {
		return time.Time{}, fmt.Errorf("invalid date format: expected 8 characters (YYYYMMDD), got %d", len(dateStr))
	}

	parsedDate, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date '%s': %w", dateStr, err)
	}

	return parsedDate, nil
}


func (h *ForecastHandler) extractRequestContext(c fiber.Ctx) (service.RequestContext, error) {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		if uid, ok := c.Locals("user_id").(int); ok {
			userID = int64(uid)
		} else {
			return service.RequestContext{}, errors.New("user_id not found or invalid type")
		}
	}

	if userID <= 0 {
		return service.RequestContext{}, errors.New("invalid user_id")
	}

	var departmentID int64
	if deptID, ok := c.Locals("department_id").(int64); ok {
		departmentID = deptID
	} else if deptID, ok := c.Locals("department_id").(int); ok {
		departmentID = int64(deptID)
	}

	return service.RequestContext{
		UserID:       userID,
		DepartmentID: departmentID,
		IPAddress:    c.IP(),
	}, nil
}

func (h *ForecastHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	forecast := router.Group("/forecasts")
	if len(ms) > 0 {
		for _, m := range ms {
			forecast.Use(m)
		}
	}
	forecast.Get("", h.GetForecast)        
	forecast.Get("/export", h.ExportReport) 
}
