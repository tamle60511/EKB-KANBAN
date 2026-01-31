package service

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type RequestContext struct {
	UserID        int64
	DepartmentID  int64
	IPAddress     string
	OperationType int
}

const (
	OperationTypeView   = 1
	OperationTypeExport = 2

	DefaultDateRangeDays = 7
	MaxDateRangeDays     = 365
)

var (
	ErrInvalidReportID   = errors.New("invalid report id")
	ErrInvalidDateRange  = errors.New("invalid date range")
	ErrDateRangeTooLarge = errors.New("date range exceeds maximum allowed days")
	ErrReportNotFound    = errors.New("report not found")
	ErrNoColumnsFound    = errors.New("no columns found for report")
	ErrEmptyReportName   = errors.New("report name cannot be empty")
	ErrEmptySQLQuery     = errors.New("sql query cannot be empty")
)

type (
	reportService struct {
		baseErpRepo   repository.BaseERPRepository
		reportRepo    repository.ReportRepo
		operationRepo repository.OperationRepository
		logger        Logger
	}
	ReportService interface {
		CreateReport(ctx context.Context, req dto.ReportCreateReq) error
		UpdateReport(ctx context.Context, id int64, req dto.ReportUpdateModel) error
		DeleteReport(ctx context.Context, id int64) error
		GetReport(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) (*dto.ReportRes, error)
		GetReportByID(ctx context.Context, id int64) (*dto.ReportDetail, error)
		GetAllReport(ctx context.Context) ([]*dto.ReportDetail, error)
		Count(ctx context.Context) (int64, error)
		ExportReport(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) (*dto.ReportFileResponse, error)
	}
)

func NewReportService(reportRepo repository.ReportRepo,
	baseErpRepo repository.BaseERPRepository, operationRepo repository.OperationRepository, logger Logger,
) ReportService {
	return &reportService{
		reportRepo:    reportRepo,
		baseErpRepo:   baseErpRepo,
		operationRepo: operationRepo,
		logger:        logger,
	}
}
func (r *reportService) CreateReport(ctx context.Context, req dto.ReportCreateReq) error {
	return r.reportRepo.CreateReport(ctx, dto.ReportCreateModel(req))
}

func (r *reportService) GetReportByID(ctx context.Context, id int64) (*dto.ReportDetail, error) {
	report, err := r.reportRepo.GetReportByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get report by id failed")
	}
	reportcolumn, err := r.reportRepo.GetColumn(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get report column error")
	}
	columns := make([]*dto.ReportColumn, len(reportcolumn))
	for i, col := range reportcolumn {
		columns[i] = &dto.ReportColumn{
			ID:       col.ID,
			ReportID: col.ReportID,
			Title:    col.Title,
			Code:     col.Code,
			Type:     col.Type,
			Num:      col.Num,
		}
	}
	return &dto.ReportDetail{
		ID:             report.ID,
		ReportType:     report.ReportType,
		ReportName:     report.ReportName,
		QueryStatement: report.QueryStatement,
		DepartmentID:   report.DepartmentID,
		Columns:        columns,
	}, nil
}
func (r *reportService) GetAllReport(ctx context.Context) ([]*dto.ReportDetail, error) {
	reports, err := r.reportRepo.GetAllReport(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all report failed")
	}
	var reportDetails []*dto.ReportDetail
	for _, report := range reports {
		reportcolumn, err := r.reportRepo.GetColumn(ctx, report.ID)
		if err != nil {
			return nil, fmt.Errorf("get report column error")
		}
		columns := make([]*dto.ReportColumn, len(reportcolumn))
		for i, col := range reportcolumn {
			columns[i] = &dto.ReportColumn{
				ID:       col.ID,
				ReportID: col.ReportID,
				Title:    col.Title,
				Code:     col.Code,
				Type:     col.Type,
				Num:      col.Num,
			}
		}
		reportDetails = append(reportDetails, &dto.ReportDetail{
			ID:             report.ID,
			ReportType:     report.ReportType,
			ReportName:     report.ReportName,
			QueryStatement: report.QueryStatement,
			DepartmentID:   report.DepartmentID,
			Columns:        columns,
		})
	}
	return reportDetails, nil
}
func (r *reportService) UpdateReport(ctx context.Context, id int64, req dto.ReportUpdateModel) error {
	return r.reportRepo.UpdateReport(ctx, id, req)
}
func (r *reportService) DeleteReport(ctx context.Context, id int64) error {
	return r.reportRepo.DeleteReport(ctx, id)
}
func (s *reportService) GetReport(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) (*dto.ReportRes, error) {
	if err := s.validateReportRequest(req); err != nil {
		return nil, err
	}

	if err := s.normalizeDateRange(req); err != nil {
		return nil, err
	}

	logID, err := s.logAccess(ctx, reqCtx, req, OperationTypeView, c)
	if err != nil {
		s.logger.Warn(ctx, "Failed to log access", map[string]interface{}{
			"error":     err.Error(),
			"user_id":   reqCtx.UserID,
			"report_id": req.ReportID,
		})
	}
	reportData, err := s.fetchReportData(ctx, req)
	if err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		return nil, err
	}
	s.updateLogStatus(ctx, logID, "success")

	return reportData, nil
}

func (s *reportService) ExportReport(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, c fiber.Ctx) (*dto.ReportFileResponse, error) {
	if err := s.validateReportRequest(req); err != nil {
		return nil, err
	}

	if err := s.normalizeDateRange(req); err != nil {
		return nil, err
	}

	logID, err := s.logAccess(ctx, reqCtx, req, OperationTypeExport, c)
	if err != nil {
		s.logger.Warn(ctx, "Failed to log access", map[string]interface{}{
			"error":     err.Error(),
			"user_id":   reqCtx.UserID,
			"report_id": req.ReportID,
		})
	}

	reportData, err := s.fetchReportData(ctx, req)
	if err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		return nil, err
	}

	reportMeta, _, err := s.getReportMetadata(ctx, req.ReportID)
	if err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		return nil, err
	}

	fileBytes, err := s.reportRepo.ExportReportToExcel(
		ctx,
		reportMeta.ReportName,
		reportData.Columns,
		reportData.Data,
		*req.FromDate,
		*req.ToDate,
	)
	if err != nil {
		s.updateLogStatus(ctx, logID, "failed")
		s.logger.Error(ctx, "Failed to generate excel file", err, map[string]interface{}{
			"report_id": req.ReportID,
		})
		return nil, fmt.Errorf("failed to generate excel file: %w", err)
	}

	fileName := s.generateFileName(reportMeta.ReportName, req.FromDate, req.ToDate)

	s.updateLogStatus(ctx, logID, "success")

	return &dto.ReportFileResponse{
		FileName:  fileName,
		FileDetal: fileBytes,
	}, nil
}
func (s *reportService) Count(ctx context.Context) (int64, error) {
	return s.reportRepo.Count(ctx)
}
func (s *reportService) fetchReportData(ctx context.Context, req *dto.ReportReq) (*dto.ReportRes, error) {
	report, columns, err := s.getReportMetadata(ctx, req.ReportID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(report.QueryStatement) == "" {
		return &dto.ReportRes{
			ReportID:   req.ReportID,
			ReportType: report.ReportType,
			ReportName: report.ReportName,
			Columns:    columns,
			Data:       nil,
		}, nil
	}

	reportDetail, err := s.baseErpRepo.GetBaseERP(ctx, dto.BaseERPReq{
		SqlQuery: report.QueryStatement,
		FromDate: *req.FromDate,
		ToDate:   *req.ToDate,
	})
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch ERP data", err, map[string]interface{}{
			"report_id": req.ReportID,
		})
		return nil, fmt.Errorf("failed to fetch report data: %w", err)
	}

	return &dto.ReportRes{
		ReportID:   req.ReportID,
		ReportType: report.ReportType,
		ReportName: report.ReportName,
		Columns:    columns,
		Data:       reportDetail.Data,
	}, nil
}

func (s *reportService) getReportMetadata(ctx context.Context, reportID int64) (*models.Report, []models.ReportColumn, error) {
	report, err := s.reportRepo.GetReport(ctx, reportID)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			return nil, nil, ErrReportNotFound
		}
		return nil, nil, fmt.Errorf("failed to get report: %w", err)
	}

	columns, err := s.reportRepo.GetColumn(ctx, reportID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get report columns: %w", err)
	}

	if len(columns) == 0 {
		return nil, nil, ErrNoColumnsFound
	}

	return report, columns, nil
}

func (s *reportService) logAccess(ctx context.Context, reqCtx RequestContext, req *dto.ReportReq, operationType int, c fiber.Ctx) (int, error) {

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

func (s *reportService) updateLogStatus(ctx context.Context, logID int, status string) {
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

func (s *reportService) normalizeDateRange(req *dto.ReportReq) error {
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

func (s *reportService) generateFileName(reportName string, fromDate, toDate *time.Time) string {
	safeName := strings.ReplaceAll(reportName, " ", "_")
	safeName = strings.ReplaceAll(safeName, "/", "-")

	return fmt.Sprintf("%s_%s_to_%s.xlsx",
		safeName,
		fromDate.Format("2006-01-02"),
		toDate.Format("2006-01-02"),
	)
}

func (s *reportService) validateReportRequest(req *dto.ReportReq) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.ReportID <= 0 {
		return ErrInvalidReportID
	}

	return nil
}
