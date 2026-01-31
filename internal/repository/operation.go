package repository

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Constants
const (
	MaxBatchSize = 1000
)

// Custom errors
var (
	ErrOperationNotFound = errors.New("operation not found")
	ErrAccessLogNotFound = errors.New("access log not found")
	ErrEmptyCode         = errors.New("code cannot be empty")
)

type (
	operationRepo struct {
		db *gorm.DB
	}

	OperationRepository interface {
		GetAll(ctx context.Context) ([]*dto.OperationResponse, error)
		FindByCode(ctx context.Context, code string) (*models.Operation, error)
		GetByID(ctx context.Context, id int) (*models.Operation, error)
		Create(ctx context.Context, operation *models.Operation) error
		Update(ctx context.Context, operation *models.Operation) error
		Delete(ctx context.Context, id int) error
		LogAccess(ctx context.Context, log *models.AccessLog) (int, error)
		UpdateLogStatus(ctx context.Context, logID int, status string) (bool, error)
		GetRecentLogs(ctx context.Context, limit int) ([]*models.AccessLog, error)
		GetLogsByUser(ctx context.Context, userID int64, fromDate, toDate time.Time, limit int) ([]*models.AccessLog, error)
		GetLogsByOperation(ctx context.Context, operationID int, fromDate, toDate time.Time, limit int) ([]*models.AccessLog, error)
		GetLogsByStatus(ctx context.Context, status string, limit int) ([]*models.AccessLog, error)
		GetLogByID(ctx context.Context, logID int) (*models.AccessLog, error)
		DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)
		LogAccessBatch(ctx context.Context, logs []*models.AccessLog) error
		GetLogCountByUser(ctx context.Context, userID int64, fromDate, toDate time.Time) (int64, error)
		GetLogCountByOperation(ctx context.Context, operationID int, fromDate, toDate time.Time) (int64, error)
		CountLogsByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
		CountFailedLogsByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
		GetAccessTrendByDays(ctx context.Context, days int, operationID *int) ([]*dto.DailyAccessCount, error)
	}
)

func NewOperationRepo(db *gorm.DB) OperationRepository {
	return &operationRepo{
		db: db,
	}
}

func (r *operationRepo) GetAll(ctx context.Context) ([]*dto.OperationResponse, error) {
	var operations []models.Operation

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&operations).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get all operations: %w", err)
	}

	// Pre-allocate slice for efficiency
	result := make([]*dto.OperationResponse, 0, len(operations))
	for _, op := range operations {
		result = append(result, &dto.OperationResponse{
			ID:   op.ID,
			Code: op.Code,
			Name: op.Name,
			Desc: op.Description,
		})
	}

	return result, nil
}

func (r *operationRepo) FindByCode(ctx context.Context, code string) (*models.Operation, error) {

	if strings.TrimSpace(code) == "" {
		return nil, ErrEmptyCode
	}

	var operation models.Operation
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&operation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOperationNotFound
		}
		return nil, fmt.Errorf("failed to find operation by code '%s': %w", code, err)
	}

	return &operation, nil
}

// GetByID retrieves an operation by ID
func (r *operationRepo) GetByID(ctx context.Context, id int) (*models.Operation, error) {
	// Validate input
	if id <= 0 {
		return nil, ErrInvalidID
	}

	var operation models.Operation
	err := r.db.WithContext(ctx).First(&operation, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOperationNotFound
		}
		return nil, fmt.Errorf("failed to get operation by id %d: %w", id, err)
	}

	return &operation, nil
}

// Create creates a new operation
func (r *operationRepo) Create(ctx context.Context, operation *models.Operation) error {
	if operation == nil {
		return errors.New("operation cannot be nil")
	}

	if err := r.db.WithContext(ctx).Create(operation).Error; err != nil {
		return fmt.Errorf("failed to create operation: %w", err)
	}

	return nil
}

// Update updates an existing operation
func (r *operationRepo) Update(ctx context.Context, operation *models.Operation) error {
	if operation == nil {
		return errors.New("operation cannot be nil")
	}

	result := r.db.WithContext(ctx).Save(operation)
	if result.Error != nil {
		return fmt.Errorf("failed to update operation: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOperationNotFound
	}

	return nil
}

// Delete deletes an operation by ID
func (r *operationRepo) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	result := r.db.WithContext(ctx).Delete(&models.Operation{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete operation: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOperationNotFound
	}

	return nil
}

// ============================================================================
// Access Log Methods
// ============================================================================

// LogAccess creates a new access log entry
func (r *operationRepo) LogAccess(ctx context.Context, log *models.AccessLog) (int, error) {
	if log == nil {
		return 0, errors.New("access log cannot be nil")
	}

	// Validate required fields
	if log.UserID <= 0 {
		return 0, errors.New("user_id is required")
	}
	if log.DepartmentID <= 0 {
		return 0, errors.New("department_id is required")
	}
	if log.OperationID <= 0 {
		return 0, errors.New("operation_id is required")
	}

	err := r.db.WithContext(ctx).Create(log).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create access log: %w", err)
	}

	return int(log.ID), nil
}

// LogAccessBatch creates multiple access log entries in a transaction
func (r *operationRepo) LogAccessBatch(ctx context.Context, logs []*models.AccessLog) error {
	if len(logs) == 0 {
		return nil
	}

	if len(logs) > MaxBatchSize {
		return fmt.Errorf("batch size %d exceeds maximum %d", len(logs), MaxBatchSize)
	}

	// Use transaction for batch insert
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(logs, 100).Error; err != nil {
			return fmt.Errorf("failed to create access logs in batch: %w", err)
		}
		return nil
	})
}

func (r *operationRepo) UpdateLogStatus(ctx context.Context, logID int, status string) (bool, error) {
	// Validate input
	if logID <= 0 {
		return false, ErrInvalidID
	}

	if strings.TrimSpace(status) == "" {
		return false, errors.New("status cannot be empty")
	}

	result := r.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Where("id = ?", logID).
		Update("status", status)

	if result.Error != nil {
		return false, fmt.Errorf("failed to update log status: %w", result.Error)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return false, ErrAccessLogNotFound
	}

	return true, nil
}

// GetRecentLogs retrieves the most recent access logs with preloaded relations
func (r *operationRepo) GetRecentLogs(ctx context.Context, limit int) ([]*models.AccessLog, error) {
	// Validate and normalize limit
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	var logs []*models.AccessLog
	err := r.db.WithContext(ctx).
		Preload("User").       // Preload user details
		Preload("Operation").  // Preload operation details
		Preload("Department"). // Preload department details (if exists)
		Order("access_time DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recent logs: %w", err)
	}

	return logs, nil
}

// GetLogsByUser retrieves access logs for a specific user
func (r *operationRepo) GetLogsByUser(
	ctx context.Context,
	userID int64,
	fromDate, toDate time.Time,
	limit int,
) ([]*models.AccessLog, error) {
	// Validate input
	if userID <= 0 {
		return nil, ErrInvalidID
	}

	if fromDate.After(toDate) {
		return nil, ErrInvalidDateRange
	}

	if limit <= 0 || limit > MaxPageSize {
		limit = MaxPageSize
	}

	var logs []*models.AccessLog
	err := r.db.WithContext(ctx).
		Preload("Operation").
		Where("user_id = ?", userID).
		Where("access_time BETWEEN ? AND ?", fromDate, toDate).
		Order("access_time DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get logs by user: %w", err)
	}

	return logs, nil
}

// GetLogsByOperation retrieves logs for a specific operation
func (r *operationRepo) GetLogsByOperation(
	ctx context.Context,
	operationID int,
	fromDate, toDate time.Time,
	limit int,
) ([]*models.AccessLog, error) {
	// Validate input
	if operationID <= 0 {
		return nil, ErrInvalidID
	}

	if fromDate.After(toDate) {
		return nil, ErrInvalidDateRange
	}

	if limit <= 0 || limit > MaxPageSize {
		limit = MaxPageSize
	}

	var logs []*models.AccessLog
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("operation_id = ?", operationID).
		Where("access_time BETWEEN ? AND ?", fromDate, toDate).
		Order("access_time DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get logs by operation: %w", err)
	}

	return logs, nil
}

// GetLogsByStatus retrieves logs by status
func (r *operationRepo) GetLogsByStatus(ctx context.Context, status string, limit int) ([]*models.AccessLog, error) {
	// Validate input
	if strings.TrimSpace(status) == "" {
		return nil, errors.New("status cannot be empty")
	}

	if limit <= 0 || limit > MaxPageSize {
		limit = MaxPageSize
	}

	var logs []*models.AccessLog
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Operation").
		Where("status = ?", status).
		Order("access_time DESC").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get logs by status: %w", err)
	}

	return logs, nil
}

// GetLogByID retrieves a single access log by ID with preloaded relations
func (r *operationRepo) GetLogByID(ctx context.Context, logID int) (*models.AccessLog, error) {
	// Validate input
	if logID <= 0 {
		return nil, ErrInvalidID
	}

	var log models.AccessLog
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Operation").
		Preload("Department").
		First(&log, logID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccessLogNotFound
		}
		return nil, fmt.Errorf("failed to get log by id: %w", err)
	}

	return &log, nil
}

// DeleteOldLogs deletes logs older than the specified date
func (r *operationRepo) DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("access_time < ?", beforeDate).
		Delete(&models.AccessLog{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old logs: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// ============================================================================
// Statistics Methods
// ============================================================================

// GetLogCountByUser counts access logs for a user in a date range
func (r *operationRepo) GetLogCountByUser(
	ctx context.Context,
	userID int64,
	fromDate, toDate time.Time,
) (int64, error) {
	// Validate input
	if userID <= 0 {
		return 0, ErrInvalidID
	}

	if fromDate.After(toDate) {
		return 0, ErrInvalidDateRange
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Where("user_id = ?", userID).
		Where("access_time BETWEEN ? AND ?", fromDate, toDate).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count logs by user: %w", err)
	}

	return count, nil
}

// GetLogCountByOperation counts access logs for an operation in a date range
func (r *operationRepo) GetLogCountByOperation(
	ctx context.Context,
	operationID int,
	fromDate, toDate time.Time,
) (int64, error) {
	// Validate input
	if operationID <= 0 {
		return 0, ErrInvalidID
	}

	if fromDate.After(toDate) {
		return 0, ErrInvalidDateRange
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Where("operation_id = ?", operationID).
		Where("access_time BETWEEN ? AND ?", fromDate, toDate).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count logs by operation: %w", err)
	}

	return count, nil
}

func (r *operationRepo) GetAccessTrendByDays(ctx context.Context, days int, operationID *int) ([]*dto.DailyAccessCount, error) {
	var results []*dto.DailyAccessCount

	startDate := time.Now().AddDate(0, 0, -days)

	query := r.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Select("DATE(access_time)::text AS date, COUNT(*)::bigint AS count").
		Where("access_time >= ?", startDate)

	// ✅ Filter by operation_id if provided
	if operationID != nil {
		query = query.Where("operation_id = ?", *operationID)
	}

	err := query.
		Group("DATE(access_time)").
		Order("DATE(access_time) ASC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get access trend: %w", err)
	}

	return results, nil
}
func (r *operationRepo) CountLogsByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Where("access_time >= ? AND access_time <= ?", startDate, endDate).
		Count(&count).Error

	return count, err
}

func (r *operationRepo) CountFailedLogsByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Where("access_time >= ? AND access_time <= ?", startDate, endDate).
		Where("status = ?", "failed").
		Count(&count).Error

	return count, err
}
