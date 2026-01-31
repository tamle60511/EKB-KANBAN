package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
	"errors"
	"fmt"
	"time"
)

type AdminService interface {
	// Operation Management
	CreateOperation(ctx context.Context, req dto.OperationCreateReq) error
	UpdateOperation(ctx context.Context, id int, req dto.OperationUpdateReq) error
	DeleteOperation(ctx context.Context, id int) error
	GetOperationByID(ctx context.Context, id int) (*dto.OperationResponse, error)
	GetAllOperationsWithStats(ctx context.Context) ([]*dto.OperationResponse, error)

	// Dashboard Analytics
	GetDashboardStats(ctx context.Context, days int) (*dto.DashboardStatsResponse, error)
	GetTopOperations(ctx context.Context, limit int, fromDate, toDate time.Time) ([]*dto.TopOperationStat, error)
	GetTopUsers(ctx context.Context, limit int, fromDate, toDate time.Time) ([]*dto.TopUserStat, error)
	GetAccessTrend(ctx context.Context, days int, operationID *int) ([]*dto.AccessTrendData, error)

	// Access Log Management
	GetAccessLogs(ctx context.Context, req dto.AccessLogQueryReq) (*dto.AccessLogListResponse, error)
	GetAccessLogByID(ctx context.Context, logID int) (*dto.AccessLogResponse, error)
	DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)

	// User Activity
	GetUserActivityReport(ctx context.Context, req dto.UserActivityReportReq) (*dto.UserActivityReportResponse, error)

	// Security
	GetSecurityAlerts(ctx context.Context, hours int) ([]*dto.SecurityAlertResponse, error)
	GetFailedAccessByIP(ctx context.Context, ipAddress string, hours int) (int64, error)
}

type adminService struct {
	operationRepo  repository.OperationRepository
	userRepo       repository.UserRepo
	departmentRepo repository.DepartmentRepo
	reportRepo     repository.ReportRepo
	logger         Logger
}

func NewAdminService(
	operationRepo repository.OperationRepository,
	userRepo repository.UserRepo,
	departmentRepo repository.DepartmentRepo,
	reportRepo repository.ReportRepo,
	logger Logger,
) AdminService {
	return &adminService{
		operationRepo:  operationRepo,
		userRepo:       userRepo,
		departmentRepo: departmentRepo,
		reportRepo:     reportRepo,
		logger:         logger,
	}
}

// ============================================================================
// Operation Management
// ============================================================================

func (s *adminService) CreateOperation(ctx context.Context, req dto.OperationCreateReq) error {
	// Check if code already exists
	existing, err := s.operationRepo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return errors.New("operation code already exists")
	}

	operation := &models.Operation{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.operationRepo.Create(ctx, operation); err != nil {
		s.logger.Error(ctx, "Failed to create operation", err, map[string]interface{}{
			"code": req.Code,
		})
		return fmt.Errorf("failed to create operation: %w", err)
	}

	s.logger.Info(ctx, "Operation created", map[string]interface{}{
		"id":   operation.ID,
		"code": operation.Code,
	})

	return nil
}

func (s *adminService) UpdateOperation(ctx context.Context, id int, req dto.OperationUpdateReq) error {
	operation, err := s.operationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	operation.Name = req.Name
	operation.Description = req.Description
	operation.UpdatedAt = time.Now()

	if err := s.operationRepo.Update(ctx, operation); err != nil {
		s.logger.Error(ctx, "Failed to update operation", err, map[string]interface{}{
			"id": id,
		})
		return fmt.Errorf("failed to update operation: %w", err)
	}

	s.logger.Info(ctx, "Operation updated", map[string]interface{}{
		"id": id,
	})

	return nil
}

func (s *adminService) DeleteOperation(ctx context.Context, id int) error {
	if err := s.operationRepo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete operation", err, map[string]interface{}{
			"id": id,
		})
		return fmt.Errorf("failed to delete operation: %w", err)
	}

	s.logger.Info(ctx, "Operation deleted", map[string]interface{}{
		"id": id,
	})

	return nil
}

func (s *adminService) GetOperationByID(ctx context.Context, id int) (*dto.OperationResponse, error) {
	operation, err := s.operationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get access count for this operation
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	accessCount, _ := s.operationRepo.GetLogCountByOperation(ctx, operation.ID, startOfMonth, now)

	return &dto.OperationResponse{
		ID:          operation.ID,
		Name:        operation.Name,
		Code:        operation.Code,
		Desc:        operation.Description,
		CreatedAt:   operation.CreatedAt,
		UpdatedAt:   operation.UpdatedAt,
		AccessCount: accessCount,
	}, nil
}

func (s *adminService) GetAllOperationsWithStats(ctx context.Context) ([]*dto.OperationResponse, error) {
	operations, err := s.operationRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Get access counts for all operations
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	result := make([]*dto.OperationResponse, 0, len(operations))
	for _, op := range operations {
		accessCount, _ := s.operationRepo.GetLogCountByOperation(ctx, op.ID, startOfMonth, now)
		op.AccessCount = accessCount
		result = append(result, op)
	}

	return result, nil
}

// ============================================================================
// Dashboard Analytics
// ============================================================================

func (s *adminService) GetDashboardStats(ctx context.Context, days int) (*dto.DashboardStatsResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// ✅ Get operations count
	operations, err := s.operationRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get operations: %w", err)
	}
	totalOperations := int64(len(operations))

	// ✅ Get counts with error handling
	totalReports, err := s.reportRepo.Count(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to count reports", err, nil)
		totalReports = 0 // Set default on error
	}

	totalUsers, err := s.userRepo.Count(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to count users", err, nil)
		totalUsers = 0
	}

	totalDepartments, err := s.departmentRepo.Count(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to count departments", err, nil)
		totalDepartments = 0
	}

	// ✅ Today's access with proper query
	todayAccess, err := s.operationRepo.CountLogsByDateRange(ctx, startOfDay, now)
	if err != nil {
		s.logger.Error(ctx, "Failed to count today access", err, nil)
		todayAccess = 0
	}

	failedAttempts, err := s.operationRepo.CountFailedLogsByDateRange(ctx, startOfDay, now)
	if err != nil {
		s.logger.Error(ctx, "Failed to count failed attempts", err, nil)
		failedAttempts = 0
	}

	topOperations, err := s.GetTopOperations(ctx, 5, startDate, now)
	if err != nil {
		s.logger.Error(ctx, "Failed to get top operations", err, nil)
		topOperations = []*dto.TopOperationStat{}
	}

	topUsers, err := s.GetTopUsers(ctx, 5, startDate, now)
	if err != nil {
		s.logger.Error(ctx, "Failed to get top users", err, nil)
		topUsers = []*dto.TopUserStat{}
	}

	recentLogs, err := s.GetRecentAccessLogs(ctx, 10)
	if err != nil {
		s.logger.Error(ctx, "Failed to get recent logs", err, nil)
		recentLogs = []*dto.AccessLogResponse{}
	}

	accessTrend, err := s.GetAccessTrend(ctx, days, nil)
	if err != nil {
		s.logger.Error(ctx, "Failed to get access trend", err, nil)
		accessTrend = []*dto.AccessTrendData{}
	}

	return &dto.DashboardStatsResponse{
		TotalOperations:  totalOperations,
		TotalUsers:       totalUsers,
		TotalDepartments: totalDepartments,
		TotalReports:     totalReports,
		TodayAccess:      todayAccess,
		FailedAttempts:   failedAttempts,
		TopOperations:    topOperations,
		TopUsers:         topUsers,
		RecentLogs:       recentLogs,
		AccessTrend:      accessTrend,
	}, nil
}

func (s *adminService) GetTopOperations(ctx context.Context, limit int, fromDate, toDate time.Time) ([]*dto.TopOperationStat, error) {
	// This would need a new repository method with GROUP BY
	// For now, we'll return a simplified version
	operations, _ := s.operationRepo.GetAll(ctx)

	stats := make([]*dto.TopOperationStat, 0)
	for _, op := range operations {
		count, _ := s.operationRepo.GetLogCountByOperation(ctx, op.ID, fromDate, toDate)
		if count > 0 {
			stats = append(stats, &dto.TopOperationStat{
				OperationID:   op.ID,
				OperationName: op.Name,
				OperationCode: op.Code,
				AccessCount:   count,
			})
		}
	}

	// Sort by access count (simplified, should use ORDER BY in SQL)
	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats, nil
}

func (s *adminService) GetTopUsers(ctx context.Context, limit int, fromDate, toDate time.Time) ([]*dto.TopUserStat, error) {

	users, _ := s.userRepo.GetAll(ctx)

	stats := make([]*dto.TopUserStat, 0)
	for _, user := range users {
		count, _ := s.operationRepo.GetLogCountByUser(ctx, user.ID, fromDate, toDate)
		if count > 0 {
			stats = append(stats, &dto.TopUserStat{
				UserID:      user.ID,
				Username:    user.Username,
				AccessCount: count,
			})
		}
	}

	// Sort by access count (simplified, should use ORDER BY in SQL)
	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats, nil
}

func (s *adminService) GetAccessTrend(ctx context.Context, days int, operationID *int) ([]*dto.AccessTrendData, error) {
	// Get data từ DB với filter
	dbResults, err := s.operationRepo.GetAccessTrendByDays(ctx, days, operationID)
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int64)
	for _, item := range dbResults {
		countMap[item.Date] = item.Count
	}

	result := make([]*dto.AccessTrendData, 0, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		result = append(result, &dto.AccessTrendData{
			Date:  dateStr,
			Count: countMap[dateStr],
		})
	}

	return result, nil
}

// ============================================================================
// Access Log Management
// ============================================================================

func (s *adminService) GetAccessLogs(ctx context.Context, req dto.AccessLogQueryReq) (*dto.AccessLogListResponse, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// This would need pagination in repository
	logs, err := s.operationRepo.GetRecentLogs(ctx, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	result := make([]*dto.AccessLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, s.convertToAccessLogResponse(log))
	}

	totalCount := int64(len(result))
	totalPages := int((totalCount + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &dto.AccessLogListResponse{
		Logs:       result,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *adminService) GetAccessLogByID(ctx context.Context, logID int) (*dto.AccessLogResponse, error) {
	log, err := s.operationRepo.GetLogByID(ctx, logID)
	if err != nil {
		return nil, err
	}

	return s.convertToAccessLogResponse(log), nil
}

func (s *adminService) DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error) {
	count, err := s.operationRepo.DeleteOldLogs(ctx, beforeDate)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete old logs", err, map[string]interface{}{
			"before_date": beforeDate,
		})
		return 0, err
	}

	s.logger.Info(ctx, "Old logs deleted", map[string]interface{}{
		"count":       count,
		"before_date": beforeDate,
	})

	return count, nil
}

func (s *adminService) GetRecentAccessLogs(ctx context.Context, limit int) ([]*dto.AccessLogResponse, error) {
	logs, err := s.operationRepo.GetRecentLogs(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.AccessLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, s.convertToAccessLogResponse(log))
	}

	return result, nil
}

// ============================================================================
// User Activity
// ============================================================================

func (s *adminService) GetUserActivityReport(ctx context.Context, req dto.UserActivityReportReq) (*dto.UserActivityReportResponse, error) {
	logs, err := s.operationRepo.GetLogsByUser(ctx, req.UserID, req.FromDate, req.ToDate, 10000)
	if err != nil {
		return nil, err
	}

	// Get user info
	user, _ := s.userRepo.GetByID(ctx, req.UserID)
	username := ""
	if user != nil {
		username = user.Username
	}

	// Calculate statistics
	var successCount, failedCount int64
	operationMap := make(map[int]int64)
	dailyMap := make(map[string]*dto.DailyActivityStat)

	for _, log := range logs {
		// Count by status
		if log.Status == StatusSuccess {
			successCount++
		} else if log.Status == StatusFailed {
			failedCount++
		}

		// Count by operation
		operationMap[log.OperationID]++

		// Count by day
		dateStr := log.AccessTime.Format("2006-01-02")
		if _, exists := dailyMap[dateStr]; !exists {
			dailyMap[dateStr] = &dto.DailyActivityStat{
				Date: dateStr,
			}
		}
		dailyMap[dateStr].AccessCount++
		if log.Status == StatusSuccess {
			dailyMap[dateStr].SuccessCount++
		} else if log.Status == StatusFailed {
			dailyMap[dateStr].FailedCount++
		}
	}

	// Convert maps to slices
	dailyActivity := make([]*dto.DailyActivityStat, 0, len(dailyMap))
	for _, stat := range dailyMap {
		dailyActivity = append(dailyActivity, stat)
	}

	operationBreakdown := make([]*dto.OperationBreakdownStat, 0, len(operationMap))
	for opID, count := range operationMap {
		op, _ := s.operationRepo.GetByID(ctx, opID)
		opName := ""
		if op != nil {
			opName = op.Name
		}
		operationBreakdown = append(operationBreakdown, &dto.OperationBreakdownStat{
			OperationID:   opID,
			OperationName: opName,
			AccessCount:   count,
		})
	}

	return &dto.UserActivityReportResponse{
		UserID:             req.UserID,
		Username:           username,
		TotalAccess:        int64(len(logs)),
		SuccessCount:       successCount,
		FailedCount:        failedCount,
		UniqueOperations:   len(operationMap),
		DailyActivity:      dailyActivity,
		OperationBreakdown: operationBreakdown,
	}, nil
}

// ============================================================================
// Security
// ============================================================================

func (s *adminService) GetSecurityAlerts(ctx context.Context, hours int) ([]*dto.SecurityAlertResponse, error) {
	fromTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	failedLogs, err := s.operationRepo.GetLogsByStatus(ctx, StatusFailed, 1000)
	if err != nil {
		return nil, err
	}

	// Group by user and IP
	alertMap := make(map[string]*dto.SecurityAlertResponse)

	for _, log := range failedLogs {
		if log.AccessTime.Before(fromTime) {
			continue
		}

		key := fmt.Sprintf("%d_%s", log.UserID, log.IPAddress)
		if alert, exists := alertMap[key]; exists {
			alert.FailedCount++
			if log.AccessTime.After(alert.LastFailedAt) {
				alert.LastFailedAt = log.AccessTime
			}
		} else {
			username := ""
			if user, err := s.userRepo.GetByID(ctx, log.UserID); err == nil && user != nil {
				username = user.Username
			}

			alertMap[key] = &dto.SecurityAlertResponse{
				AlertType:    "multiple_failed_attempts",
				UserID:       log.UserID,
				Username:     username,
				IPAddress:    log.IPAddress,
				FailedCount:  1,
				LastFailedAt: log.AccessTime,
				Description:  fmt.Sprintf("Multiple failed access attempts from IP %s", log.IPAddress),
			}
		}
	}

	// Convert to slice and filter (only alerts with 3+ failures)
	alerts := make([]*dto.SecurityAlertResponse, 0)
	for _, alert := range alertMap {
		if alert.FailedCount >= 3 {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

func (s *adminService) GetFailedAccessByIP(ctx context.Context, ipAddress string, hours int) (int64, error) {
	fromTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	failedLogs, err := s.operationRepo.GetLogsByStatus(ctx, StatusFailed, 10000)
	if err != nil {
		return 0, err
	}

	count := int64(0)
	for _, log := range failedLogs {
		if log.IPAddress == ipAddress && log.AccessTime.After(fromTime) {
			count++
		}
	}

	return count, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

func (s *adminService) convertToAccessLogResponse(log *models.AccessLog) *dto.AccessLogResponse {
	resp := &dto.AccessLogResponse{
		ID:           log.ID,
		UserID:       log.UserID,
		DepartmentID: log.DepartmentID,
		OperationID:  log.OperationID,
		AccessTime:   log.AccessTime,
		ReportID:     log.ReportID,
		IPAddress:    log.IPAddress,
		Status:       log.Status,
		CreatedAt:    log.CreatedAt,
	}

	if user, err := s.userRepo.GetByID(context.Background(), log.UserID); err == nil && user != nil {
		resp.Username = user.Username
	}

	if operation, err := s.operationRepo.GetByID(context.Background(), log.OperationID); err == nil && operation != nil {
		resp.OperationName = operation.Name
		resp.OperationCode = operation.Code
	}

	return resp
}
