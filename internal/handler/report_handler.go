package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type ReportHandler struct {
	BaseHandler
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (r *ReportHandler) CreateReport(c fiber.Ctx) error {
	var report dto.ReportCreateReq
	if err := c.Bind().Body(&report); err != nil {
		return utils.BadRequestResponse(c, "invalid body parser", err.Error())
	}
	err := r.reportService.CreateReport(c.RequestCtx(), report)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to create report", err)
	}
	return utils.SuccessResponse(c, "create report success", nil)
}
func (r *ReportHandler) UpdateReport(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid ID params")
	}
	var req dto.ReportUpdateModel
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := r.reportService.UpdateReport(c.RequestCtx(), int64(id), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update report", err)
	}
	return utils.SuccessResponse(c, "Update to report success", nil)
}
func (r *ReportHandler) DeleteReport(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid ID params")
	}
	if err := r.reportService.DeleteReport(c.RequestCtx(), id); err != nil {
		return utils.BadRequestResponse(c, "failed to delete", err)
	}
	return utils.SuccessResponse(c, "delete success", nil)
}
func (r *ReportHandler) GetReport(c fiber.Ctx) error {
	var req dto.ReportReq
	if err := c.Bind().Query(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := c.Bind().URI(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid params", err)
	}
	reqCtx, err := r.extractRequestContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}
	report, err := r.reportService.GetReport(c.RequestCtx(), reqCtx, &req, c)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed get report columns", err)
	}
	return utils.SuccessResponse(c, "get report success", report)
}

func (r *ReportHandler) ExportReport(c fiber.Ctx) error {

	var req dto.ReportReq
	if err := c.Bind().Query(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err)
	}
	if err := c.Bind().URI(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid params", err)
	}
	reqCtx, err := r.extractRequestContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}
	reportData, err := r.reportService.ExportReport(c.RequestCtx(), reqCtx, &req, c)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed export report", err)
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, reportData.FileName))
	if fileBytes, ok := reportData.FileDetal.([]byte); ok {
		c.Set("Content-Length", strconv.Itoa(len(fileBytes)))

		return c.Send(fileBytes)
	} else {
		return utils.InternalErrorResponse(c, "Invalid file data", "reportData.FileDetal is not []byte")
	}
}

func (r *ReportHandler) GetReportByID(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid ID params")
	}
	report, err := r.reportService.GetReportByID(c.RequestCtx(), id)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed get report by id", err)
	}
	return utils.SuccessResponse(c, "get report success", report)
}
func (r *ReportHandler) GetAllReport(c fiber.Ctx) error {
	reports, err := r.reportService.GetAllReport(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed get all report", err)
	}
	return utils.SuccessResponse(c, "get all report success", reports)
}

func (h *ReportHandler) extractRequestContext(c fiber.Ctx) (service.RequestContext, error) {
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

func (r *ReportHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	report := router.Group("/reports")
	if len(ms) > 0 {
		for _, m := range ms {
			report.Use(m)
		}
	}
	report.Post("/", r.CreateReport)
	report.Get("/", r.GetAllReport)
	report.Get("/:id", r.GetReport)
	report.Put("/:id", r.UpdateReport)
	report.Get("/by-id/:id", r.GetReportByID)
	report.Delete("/:id", r.DeleteReport)
	report.Get("/export/:id", r.ExportReport)
}
