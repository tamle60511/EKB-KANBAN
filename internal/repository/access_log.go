package repository

import (
	"context"
	"cqs-kanban/internal/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Constants for pagination and batch operations
const (
	DefaultPageSize  = 10
	MaxPageSize      = 100
	DefaultBatchSize = 100
	DefaultLimit     = 10
)

// Custom errors
var (
	ErrInvalidDateRange = errors.New("fromDate must be before toDate")
	ErrInvalidLimit     = errors.New("limit must be positive")
	ErrInvalidPageSize  = errors.New("page size must be positive")
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidID        = errors.New("invalid id")
	ErrLogCannotBeNil   = errors.New("log cannot be nil")
	ErrStatusEmpty      = errors.New("status cannot be empty")
)

// PaginationParams defines pagination parameters
type PaginationParams struct {
	Page     int // Current page (starts from 1)
	PageSize int // Number of items per page
}

// Validate validates pagination parameters
func (p *PaginationParams) Validate() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}

// PaginatedResult represents paginated result
type PaginatedResult struct {
	Data       []models.AccessLog `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

type (
	accessLogRepo struct {
		db *gorm.DB
	}

	// AccessLogRepo defines the interface for access log operations
	// Recommended database indexes:
	// - idx_user_id_access_time ON (user_id, access_time DESC)
	// - idx_user_id_status ON (user_id, status)
	// - idx_access_time ON (access_time)
	AccessLogRepo interface {
		// CreateAccessLog creates a new access log entry
		CreateAccessLog(ctx context.Context, log *models.AccessLog) error

		// CreateAccessLogsBatch creates multiple access logs in batches
		CreateAccessLogsBatch(ctx context.Context, logs []*models.AccessLog, batchSize int) error

		// UpdateStatus updates the status of an access log
		UpdateStatus(ctx context.Context, id int, status string) error

		// GetByID retrieves an access log by its ID
		GetByID(ctx context.Context, id int) (*models.AccessLog, error)

		// GetRecentLogs retrieves the most recent logs for a user
		GetRecentLogs(ctx context.Context, userID int, limit int) ([]models.AccessLog, error)

		// GetUserLogs retrieves all logs for a user within a date range
		GetUserLogs(ctx context.Context, userID int, fromDate, toDate time.Time) ([]models.AccessLog, error)

		// GetUserLogsPaginated retrieves logs for a user with pagination
		GetUserLogsPaginated(ctx context.Context, userID int, fromDate, toDate time.Time, pagination PaginationParams) (*PaginatedResult, error)

		// DeleteOldLogs deletes logs older than the specified time
		DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error)

		// GetLogsByStatus retrieves logs filtered by status
		GetLogsByStatus(ctx context.Context, userID int, status string, limit int) ([]models.AccessLog, error)
	}
)

// NewAccessLogRepo creates a new AccessLogRepo instance
func NewAccessLogRepo(db *gorm.DB) AccessLogRepo {
	return &accessLogRepo{db: db}
}

// CreateAccessLog creates a new access log entry in the database
func (a *accessLogRepo) CreateAccessLog(ctx context.Context, log *models.AccessLog) error {
	if log == nil {
		return ErrLogCannotBeNil
	}

	if err := a.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create access log: %w", err)
	}
	return nil
}

// CreateAccessLogsBatch creates multiple access logs using batch insert
// batchSize: number of records to insert per batch (0 or negative uses default)
func (a *accessLogRepo) CreateAccessLogsBatch(ctx context.Context, logs []*models.AccessLog, batchSize int) error {
	if len(logs) == 0 {
		return nil
	}

	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}

	if err := a.db.WithContext(ctx).CreateInBatches(logs, batchSize).Error; err != nil {
		return fmt.Errorf("failed to create access logs batch: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of an access log by ID
func (a *accessLogRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if status == "" {
		return ErrStatusEmpty
	}

	result := a.db.WithContext(ctx).
		Model(&models.AccessLog{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetByID retrieves an access log by its ID
func (a *accessLogRepo) GetByID(ctx context.Context, id int) (*models.AccessLog, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	var log models.AccessLog
	err := a.db.WithContext(ctx).
		Where("id = ?", id).
		First(&log).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get access log by id: %w", err)
	}

	return &log, nil
}

// GetRecentLogs retrieves the most recent access logs for a user
// limit: maximum number of logs to return (0 uses default, max 100)
func (a *accessLogRepo) GetRecentLogs(ctx context.Context, userID int, limit int) ([]models.AccessLog, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	// Apply limit constraints
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	var logs []models.AccessLog
	err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("access_time desc").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recent logs: %w", err)
	}

	return logs, nil
}

// GetUserLogs retrieves all access logs for a user within a date range
func (a *accessLogRepo) GetUserLogs(ctx context.Context, userID int, fromDate, toDate time.Time) ([]models.AccessLog, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if fromDate.After(toDate) {
		return nil, ErrInvalidDateRange
	}

	var logs []models.AccessLog
	err := a.db.WithContext(ctx).
		Where("user_id = ? AND access_time >= ? AND access_time < ?", userID, fromDate, toDate).
		Order("access_time desc").
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user logs: %w", err)
	}

	return logs, nil
}

func (a *accessLogRepo) GetUserLogsPaginated(ctx context.Context, userID int, fromDate, toDate time.Time, pagination PaginationParams) (*PaginatedResult, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if fromDate.After(toDate) {
		return nil, ErrInvalidDateRange
	}

	pagination.Validate()
	query := a.db.WithContext(ctx).
		Where("user_id = ? AND access_time >= ? AND access_time < ?", userID, fromDate, toDate)

	var total int64
	if err := query.Model(&models.AccessLog{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count user logs: %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + pagination.PageSize - 1) / pagination.PageSize
	}

	var logs []models.AccessLog
	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Order("access_time desc").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get paginated user logs: %w", err)
	}

	return &PaginatedResult{
		Data:       logs,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (a *accessLogRepo) DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	result := a.db.WithContext(ctx).
		Where("access_time < ?", olderThan).
		Delete(&models.AccessLog{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old logs: %w", result.Error)
	}

	return result.RowsAffected, nil
}

func (a *accessLogRepo) GetLogsByStatus(ctx context.Context, userID int, status string, limit int) ([]models.AccessLog, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if status == "" {
		return nil, ErrStatusEmpty
	}

	// Apply limit constraints
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	var logs []models.AccessLog
	err := a.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Order("access_time desc").
		Limit(limit).
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get logs by status: %w", err)
	}

	return logs, nil
}
