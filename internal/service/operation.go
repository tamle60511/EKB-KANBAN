package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	DefaultLogLimit = 10
	MaxLogLimit     = 100
	MaxIPLength     = 45
	MaxParamsLength = 5000
	CacheTTL        = 5 * time.Minute
)

const (
	StatusPending   = "pending"
	StatusSuccess   = "success"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

var validStatuses = map[string]bool{
	StatusPending:   true,
	StatusSuccess:   true,
	StatusFailed:    true,
	StatusCancelled: true,
}

var (
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidLogID       = errors.New("invalid log id")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidIPAddress   = errors.New("invalid ip address")
	ErrEmptyOperationCode = errors.New("operation code cannot be empty")
	ErrOperationNotFound  = errors.New("operation not found")
	ErrParamsTooLarge     = errors.New("params exceed maximum size")
)

type OperationCache struct {
	mu    sync.RWMutex
	cache map[string]*models.Operation
	ttl   time.Duration
}

func NewOperationCache(ttl time.Duration) *OperationCache {
	return &OperationCache{
		cache: make(map[string]*models.Operation),
		ttl:   ttl,
	}
}

func (c *OperationCache) Get(code string) (*models.Operation, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	op, exists := c.cache[code]
	return op, exists
}

func (c *OperationCache) Set(code string, op *models.Operation) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[code] = op
}

func (c *OperationCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*models.Operation)
}

type (
	operationService struct {
		operationRepo repository.OperationRepository
		userRepo      repository.UserRepo
		cache         *OperationCache
		logger        Logger
	}

	OperationService interface {
		GetAllOperations(ctx context.Context) ([]*dto.OperationResponse, error)
		LogAccess(ctx context.Context, userID int64, departmentID int64, operationCode string, reportID int64, c fiber.Ctx) (int, error)
		UpdateLogStatus(ctx context.Context, logID int, status string) error
		GetRecentLogs(ctx context.Context, limit int) ([]*models.AccessLog, error)
		GetLogsByUser(ctx context.Context, userID int64, fromDate, toDate time.Time, limit int) ([]*models.AccessLog, error)
		GetLogsByOperation(ctx context.Context, operationID int, fromDate, toDate time.Time, limit int) ([]*models.AccessLog, error)
		GetLogsByStatus(ctx context.Context, status string, limit int) ([]*models.AccessLog, error)
		GetOperationByCode(ctx context.Context, code string) (*models.Operation, error)
		ClearCache()
	}
)

// ✅ Fixed constructor
func NewOperationService(
	operationRepo repository.OperationRepository,
	userRepo repository.UserRepo,
	logger Logger,

) OperationService {
	return &operationService{
		operationRepo: operationRepo,
		userRepo:      userRepo,
		cache:         NewOperationCache(CacheTTL),
		logger:        logger,
	}
}

func (s *operationService) GetAllOperations(ctx context.Context) ([]*dto.OperationResponse, error) {
	operations, err := s.operationRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get all operations", err, nil)
		return nil, fmt.Errorf("failed to get operations: %w", err)
	}

	return operations, nil
}

func (s *operationService) LogAccess(
	ctx context.Context,
	userID int64,
	departmentID int64,
	operationCode string,
	reportID int64,
	c fiber.Ctx,
) (int, error) {
	if err := s.validateLogAccessInput(userID, operationCode, c.IP()); err != nil {
		return 0, err
	}

	operation, err := s.GetOperationByCode(ctx, operationCode)
	if err != nil {
		s.logger.Error(ctx, "Operation not found", err, map[string]interface{}{
			"operation_code": operationCode,
		})
		return 0, fmt.Errorf("operation lookup failed: %w", err)
	}

	ipAddress := c.IP()
	log := &models.AccessLog{
		UserID:       userID,
		DepartmentID: departmentID,
		OperationID:  operation.ID,
		AccessTime:   time.Now().UTC(),
		ReportID:     int64(reportID),
		IPAddress:    ipAddress,
		Status:       StatusPending,
	}

	logID, err := s.operationRepo.LogAccess(ctx, log)
	if err != nil {
		s.logger.Error(ctx, "Failed to log access", err, map[string]interface{}{
			"user_id":        userID,
			"operation_code": operationCode,
		})
		return 0, fmt.Errorf("failed to log access: %w", err)
	}

	s.logger.Info(ctx, "Access logged successfully", map[string]interface{}{
		"log_id":         logID,
		"user_id":        userID,
		"operation_code": operationCode,
	})

	return logID, nil
}

func (s *operationService) UpdateLogStatus(ctx context.Context, logID int, status string) error {
	if logID <= 0 {
		return ErrInvalidLogID
	}

	if !s.isValidStatus(status) {
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	updated, err := s.operationRepo.UpdateLogStatus(ctx, logID, status)
	if err != nil {
		s.logger.Error(ctx, "Failed to update log status", err, map[string]interface{}{
			"log_id": logID,
			"status": status,
		})
		return fmt.Errorf("failed to update log status: %w", err)
	}

	if !updated {
		return fmt.Errorf("log %d not found or not updated", logID)
	}

	return nil
}

func (s *operationService) GetRecentLogs(ctx context.Context, limit int) ([]*models.AccessLog, error) {
	limit = s.normalizeLimit(limit)

	logs, err := s.operationRepo.GetRecentLogs(ctx, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to get recent logs", err, map[string]interface{}{
			"limit": limit,
		})
		return nil, fmt.Errorf("failed to get recent logs: %w", err)
	}

	return logs, nil
}

func (s *operationService) GetLogsByUser(
	ctx context.Context,
	userID int64,
	fromDate, toDate time.Time,
	limit int,
) ([]*models.AccessLog, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if fromDate.After(toDate) {
		return nil, errors.New("fromDate must be before toDate")
	}

	limit = s.normalizeLimit(limit)

	logs, err := s.operationRepo.GetLogsByUser(ctx, userID, fromDate, toDate, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to get logs by user", err, map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get user logs: %w", err)
	}

	return logs, nil
}

func (s *operationService) GetLogsByOperation(
	ctx context.Context,
	operationID int,
	fromDate, toDate time.Time,
	limit int,
) ([]*models.AccessLog, error) {
	if operationID <= 0 {
		return nil, errors.New("invalid operation id")
	}

	if fromDate.After(toDate) {
		return nil, errors.New("fromDate must be before toDate")
	}

	limit = s.normalizeLimit(limit)

	logs, err := s.operationRepo.GetLogsByOperation(ctx, operationID, fromDate, toDate, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to get logs by operation", err, map[string]interface{}{
			"operation_id": operationID,
		})
		return nil, fmt.Errorf("failed to get operation logs: %w", err)
	}

	return logs, nil
}

func (s *operationService) GetLogsByStatus(ctx context.Context, status string, limit int) ([]*models.AccessLog, error) {
	if !s.isValidStatus(status) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	limit = s.normalizeLimit(limit)

	logs, err := s.operationRepo.GetLogsByStatus(ctx, status, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to get logs by status", err, map[string]interface{}{
			"status": status,
		})
		return nil, fmt.Errorf("failed to get logs by status: %w", err)
	}

	return logs, nil
}

func (s *operationService) GetOperationByCode(ctx context.Context, code string) (*models.Operation, error) {
	if strings.TrimSpace(code) == "" {
		return nil, ErrEmptyOperationCode
	}

	if cached, exists := s.cache.Get(code); exists {
		return cached, nil
	}

	operation, err := s.operationRepo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrOperationNotFound) {
			return nil, ErrOperationNotFound
		}
		return nil, fmt.Errorf("failed to find operation: %w", err)
	}

	s.cache.Set(code, operation)

	return operation, nil
}

func (s *operationService) ClearCache() {
	s.cache.Clear()
	s.logger.Info(context.Background(), "Operation cache cleared", nil)
}

func (s *operationService) validateLogAccessInput(userID int64, operationCode, ipAddress string) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}

	if strings.TrimSpace(operationCode) == "" {
		return ErrEmptyOperationCode
	}

	if !s.isValidIPAddress(ipAddress) {
		return ErrInvalidIPAddress
	}

	return nil
}

func (s *operationService) isValidIPAddress(ip string) bool {
	if ip == "" {
		return false
	}

	if len(ip) > MaxIPLength {
		return false
	}

	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func (s *operationService) isValidStatus(status string) bool {
	return validStatuses[status]
}

func (s *operationService) normalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLogLimit
	}
	if limit > MaxLogLimit {
		return MaxLogLimit
	}
	return limit
}

func (s *operationService) marshalParams(params interface{}) string {
	if params == nil {
		return ""
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		s.logger.Warn(context.Background(), "Failed to marshal params", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Sprintf(`{"error": "failed to marshal params: %s"}`, err.Error())
	}

	if len(jsonBytes) > MaxParamsLength {
		s.logger.Warn(context.Background(), "Params exceed max length", map[string]interface{}{
			"length": len(jsonBytes),
			"max":    MaxParamsLength,
		})
		return fmt.Sprintf(`{"error": "params too large", "size": %d}`, len(jsonBytes))
	}

	return string(jsonBytes)
}
