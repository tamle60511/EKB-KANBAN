package handler

import (
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/service"
	"cqs-kanban/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type AdminHandler struct {
	BaseHandler
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// getQueryInt lấy query parameter kiểu int với giá trị mặc định
func getQueryInt(c fiber.Ctx, key string, defaultValue int) int {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return result
}

// ============================================================================
// Operation Management Endpoints
// ============================================================================

func (h *AdminHandler) CreateOperation(c fiber.Ctx) error {
	var req dto.OperationCreateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.adminService.CreateOperation(c.RequestCtx(), req); err != nil {
		return utils.InternalErrorResponse(c, "Failed to create operation", err)
	}

	return utils.SuccessResponse(c, "Operation created successfully", nil)
}

func (h *AdminHandler) UpdateOperation(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid operation ID", err.Error())
	}

	var req dto.OperationUpdateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.adminService.UpdateOperation(c.RequestCtx(), id, req); err != nil {
		return utils.InternalErrorResponse(c, "Failed to update operation", err)
	}

	return utils.SuccessResponse(c, "Operation updated successfully", nil)
}

func (h *AdminHandler) DeleteOperation(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid operation ID", err.Error())
	}

	if err := h.adminService.DeleteOperation(c.RequestCtx(), id); err != nil {
		return utils.InternalErrorResponse(c, "Failed to delete operation", err)
	}

	return utils.SuccessResponse(c, "Operation deleted successfully", nil)
}

func (h *AdminHandler) GetOperation(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid operation ID", err.Error())
	}

	operation, err := h.adminService.GetOperationByID(c.RequestCtx(), id)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get operation", err)
	}

	return utils.SuccessResponse(c, "Operation retrieved successfully", operation)
}

func (h *AdminHandler) GetAllOperations(c fiber.Ctx) error {
	operations, err := h.adminService.GetAllOperationsWithStats(c.RequestCtx())
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get operations", err)
	}

	return utils.SuccessResponse(c, "Operations retrieved successfully", operations)
}

// ============================================================================
// Dashboard Analytics Endpoints
// ============================================================================

func (h *AdminHandler) GetDashboard(c fiber.Ctx) error {
	days := getQueryInt(c, "days", 7)

	stats, err := h.adminService.GetDashboardStats(c.RequestCtx(), days)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get dashboard stats", err)
	}

	return utils.SuccessResponse(c, "Dashboard stats retrieved successfully", stats)
}

func (h *AdminHandler) GetTopOperations(c fiber.Ctx) error {
	limit := getQueryInt(c, "limit", 10)
	days := getQueryInt(c, "days", 30)

	now := time.Now()
	fromDate := now.AddDate(0, 0, -days)

	stats, err := h.adminService.GetTopOperations(c.RequestCtx(), limit, fromDate, now)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get top operations", err)
	}

	return utils.SuccessResponse(c, "Top operations retrieved successfully", stats)
}

func (h *AdminHandler) GetTopUsers(c fiber.Ctx) error {
	limit := getQueryInt(c, "limit", 10)
	days := getQueryInt(c, "days", 30)

	now := time.Now()
	fromDate := now.AddDate(0, 0, -days)

	stats, err := h.adminService.GetTopUsers(c.RequestCtx(), limit, fromDate, now)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get top users", err)
	}

	return utils.SuccessResponse(c, "Top users retrieved successfully", stats)
}

func (h *AdminHandler) GetAccessTrend(c fiber.Ctx) error {
	days := getQueryInt(c, "days", 7)
	var operationID *int
	if opIDStr := c.Query("operation_id"); opIDStr != "" {
		if opID, err := strconv.Atoi(opIDStr); err == nil {
			operationID = &opID
		} else {
			return utils.BadRequestResponse(c, "Invalid operation_id parameter", err.Error())
		}
	}
	trend, err := h.adminService.GetAccessTrend(c.RequestCtx(), days, operationID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get access trend", err)
	}

	return utils.SuccessResponse(c, "Access trend retrieved successfully", trend)
}

// ============================================================================
// Access Log Endpoints
// ============================================================================

func (h *AdminHandler) GetAccessLogs(c fiber.Ctx) error {
	var req dto.AccessLogQueryReq
	if err := c.Bind().Query(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid query parameters", err.Error())
	}

	logs, err := h.adminService.GetAccessLogs(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get access logs", err)
	}

	return utils.SuccessResponse(c, "Access logs retrieved successfully", logs)
}

func (h *AdminHandler) GetAccessLog(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid log ID", err.Error())
	}

	log, err := h.adminService.GetAccessLogByID(c.RequestCtx(), id)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get access log", err)
	}

	return utils.SuccessResponse(c, "Access log retrieved successfully", log)
}

func (h *AdminHandler) DeleteOldLogs(c fiber.Ctx) error {
	days := getQueryInt(c, "days", 90)

	beforeDate := time.Now().AddDate(0, 0, -days)
	count, err := h.adminService.DeleteOldLogs(c.RequestCtx(), beforeDate)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to delete old logs", err)
	}

	return utils.SuccessResponse(c, "Old logs deleted successfully", fiber.Map{
		"deleted_count": count,
		"before_date":   beforeDate,
	})
}

// ============================================================================
// User Activity Endpoints
// ============================================================================

func (h *AdminHandler) GetUserActivityReport(c fiber.Ctx) error {
	var req dto.UserActivityReportReq
	if err := c.Bind().Query(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid query parameters", err.Error())
	}

	report, err := h.adminService.GetUserActivityReport(c.RequestCtx(), req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get user activity report", err)
	}

	return utils.SuccessResponse(c, "User activity report retrieved successfully", report)
}

// ============================================================================
// Security Endpoints
// ============================================================================

func (h *AdminHandler) GetSecurityAlerts(c fiber.Ctx) error {
	hours := getQueryInt(c, "hours", 24)

	alerts, err := h.adminService.GetSecurityAlerts(c.RequestCtx(), hours)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get security alerts", err)
	}

	return utils.SuccessResponse(c, "Security alerts retrieved successfully", alerts)
}

func (h *AdminHandler) GetFailedAccessByIP(c fiber.Ctx) error {
	ipAddress := c.Query("ip")
	if ipAddress == "" {
		return utils.BadRequestResponse(c, "IP address is required", "")
	}

	hours := getQueryInt(c, "hours", 24)

	count, err := h.adminService.GetFailedAccessByIP(c.RequestCtx(), ipAddress, hours)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get failed access count", err)
	}

	return utils.SuccessResponse(c, "Failed access count retrieved", fiber.Map{
		"ip_address":   ipAddress,
		"failed_count": count,
		"hours":        hours,
	})
}

// ============================================================================
// Setup Routes
// ============================================================================

func (h *AdminHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	admin := router.Group("/admin")

	// Apply middleware
	if len(ms) > 0 {
		for _, m := range ms {
			admin.Use(m)
		}
	}

	// Dashboard
	admin.Get("/dashboard", h.GetDashboard)
	admin.Get("/dashboard/top-operations", h.GetTopOperations)
	admin.Get("/dashboard/top-users", h.GetTopUsers)
	admin.Get("/dashboard/access-trend", h.GetAccessTrend)

	// Operations Management
	operations := admin.Group("/operations")
	operations.Post("/", h.CreateOperation)
	operations.Get("/", h.GetAllOperations)
	operations.Get("/:id", h.GetOperation)
	operations.Put("/:id", h.UpdateOperation)
	operations.Delete("/:id", h.DeleteOperation)

	// Access Logs
	logs := admin.Group("/logs")
	logs.Get("/", h.GetAccessLogs)
	logs.Get("/:id", h.GetAccessLog)
	logs.Delete("/cleanup", h.DeleteOldLogs)

	// User Activity
	admin.Get("/user-activity", h.GetUserActivityReport)

	// Security
	security := admin.Group("/security")
	security.Get("/alerts", h.GetSecurityAlerts)
	security.Get("/failed-access", h.GetFailedAccessByIP)
}
