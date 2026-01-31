package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

type (
	forecastService struct {
		forecastRepo  repository.ForecastRepo
		operationRepo repository.OperationRepository
		logger        Logger
	}
	ForecastService interface {
		GetForecast(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) ([]dto.CombinedForecast, error)
		ExportReport(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) (*dto.ReportFileResponse, error)
	}
)

func NewForecastService(forecastRepo repository.ForecastRepo, operationRepo repository.OperationRepository, logger Logger) ForecastService {
	return &forecastService{forecastRepo: forecastRepo, operationRepo: operationRepo, logger: logger}
}

func (s *forecastService) GetForecast(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) ([]dto.CombinedForecast, error) {
	logID, err := s.logAccess(ctx, reqCtx, req, OperationTypeView, c)
	if err != nil {
		s.logger.Warn(ctx, "Failed to log access", map[string]interface{}{
			"error":     err.Error(),
			"user_id":   reqCtx.UserID,
			"report_id": req.ReportID,
		})
	}
	if err := s.normalizeDateRange(req); err != nil {
		return nil, err
	}
	s.updateLogStatus(ctx, logID, "success")
	return s.forecastRepo.GetForecast(ctx, req)
}

func (s *forecastService) ExportReport(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) (*dto.ReportFileResponse, error) {
	logID, err := s.logAccess(ctx, reqCtx, req, OperationTypeExport, c)
	if err != nil {
		s.logger.Warn(ctx, "Failed to log access", map[string]interface{}{
			"error":     err.Error(),
			"user_id":   reqCtx.UserID,
			"report_id": req.ReportID,
		})
	}

	if err := s.normalizeDateRange(req); err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		return nil, err
	}
	forecastData, err := s.forecastRepo.GetForecast(ctx, req)
	if err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		return nil, fmt.Errorf("failed to fetch forecast data: %w", err)
	}

	fromDate := *req.FromDate
	toDate := *req.ToDate

	fileBytes, err := s.forecastRepo.ExportGroupedForecastToExcel(ctx, forecastData, fromDate, toDate)
	if err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		return nil, fmt.Errorf("failed to export to excel: %w", err)
	}

	s.updateLogStatus(ctx, logID, "success")

	reportName := "Forecast_Report"
	return &dto.ReportFileResponse{
		ReportName:  reportName,
		FileName:    fmt.Sprintf("%s_%s_to_%s.xlsx", reportName, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")),
		FileDetal:   fileBytes,
		GeneratedAt: time.Now(),
	}, nil
}
func (s *forecastService) normalizeDateRange(req *dto.ReportReq) error {
	now := time.Now().UTC()
	if req.FromDate == nil {
		fromDate := now.AddDate(0, 0, -DefaultDateRangeDays)
		req.FromDate = &fromDate
	}
	if req.ToDate == nil {
		req.ToDate = &now
	}

	if req.FromDate.After(*req.ToDate) {
		return ErrInvalidDateRange
	}

	daysDiff := int(req.ToDate.Sub(*req.FromDate).Hours() / 24)
	if daysDiff > MaxDateRangeDays {
		return fmt.Errorf("%w: maximum %d days allowed, got %d days", ErrDateRangeTooLarge, MaxDateRangeDays, daysDiff)
	}

	return nil
}

func (s *forecastService) logAccess(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, operationType int, c fiber.Ctx) (int, error) {
	ipAddress := c.IP()
	accessLog := &models.AccessLog{
		UserID:       reqCtx.UserID,
		DepartmentID: reqCtx.DepartmentID,
		OperationID:  operationType,
		AccessTime:   time.Now().UTC(),
		IPAddress:    ipAddress,
		ReportID:     req.ReportID,
		Status:       "pending",
	}

	logID, err := s.operationRepo.LogAccess(ctx, accessLog)
	if err != nil {
		return 0, fmt.Errorf("failed to log access: %w", err)
	}

	return logID, nil
}

func (s *forecastService) updateLogStatus(ctx context.Context, logID int, status string) {
	if logID <= 0 {
		return
	}

	if _, err := s.operationRepo.UpdateLogStatus(ctx, logID, status); err != nil {
		s.logger.Error(ctx, "Failed to update log status", err, map[string]interface{}{
			"log_id": logID,
			"status": status,
		})
	}
}


