package repository

import (
	"context"
	"cqs-kanban/internal/dto"
	"cqs-kanban/internal/models"
	"cqs-kanban/internal/utils"
	"fmt"
	"io/ioutil"
	"time"

	"gorm.io/gorm"
)

type (
	reportRepo struct {
		db *gorm.DB
	}
	ReportRepo interface {
		CreateReport(ctx context.Context, input dto.ReportCreateModel) error
		UpdateReport(ctx context.Context, id int64, input dto.ReportUpdateModel) error
		DeleteReport(ctx context.Context, id int64) error
		GetReport(ctx context.Context, id int64) (*models.Report, error)
		GetReportByID(ctx context.Context, id int64) (*models.Report, error)
		GetAllReport(ctx context.Context) ([]models.Report, error)
		GetColumn(ctx context.Context, reportID int64) ([]models.ReportColumn, error)
		Count(ctx context.Context) (int64, error)
		ExportReportToExcel(ctx context.Context, reportName string, columns []models.ReportColumn, data []map[string]interface{}, fromDate, toDate time.Time) ([]byte, error)
	}
)

func NewReportRepo(db *gorm.DB) ReportRepo {
	return &reportRepo{
		db: db,
	}
}

func (r *reportRepo) CreateReport(ctx context.Context, input dto.ReportCreateModel) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		report := models.Report{
			ReportType:     input.ReportType,
			ReportName:     input.ReportName,
			QueryStatement: input.QueryStatement,
			DepartmentID:   input.DepartmentID,
		}
		if err := tx.WithContext(ctx).Create(&report).Error; err != nil {
			return fmt.Errorf("failed to create report")
		}
		for _, column := range input.Columns {
			column.ReportID = report.ID
		}
		if err := tx.WithContext(ctx).Create(&input.Columns).Error; err != nil {
			return fmt.Errorf("failed to created column report")
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error transaction create report")
	}
	return nil
}

func (r *reportRepo) GetReport(ctx context.Context, id int64) (*models.Report, error) {
	report, err := gorm.G[models.Report](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("get report error")
	}
	return &report, nil
}
func (r *reportRepo) GetColumn(ctx context.Context, reportID int64) ([]models.ReportColumn, error) {
	column, err := gorm.G[models.ReportColumn](r.db).Where("report_id = ?", reportID).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("get column report error")
	}
	return column, nil
}
func (r *reportRepo) GetReportByID(ctx context.Context, id int64) (*models.Report, error) {
	report, err := gorm.G[models.Report](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("get report by id error")
	}
	return &report, nil
}
func (r *reportRepo) GetAllReport(ctx context.Context) ([]models.Report, error) {
	reports, err := gorm.G[models.Report](r.db).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all report error")
	}
	return reports, nil
}
func (r *reportRepo) UpdateReport(ctx context.Context, id int64, input dto.ReportUpdateModel) error {
	report, err := r.GetReportByID(ctx, id)
	if err != nil {
		return fmt.Errorf("report not found")
	}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		report.ReportType = *input.ReportType
		report.ReportName = *input.ReportName
		report.QueryStatement = *input.QueryStatement
		report.DepartmentID = *input.DepartmentID
		if err := tx.WithContext(ctx).Save(&report).Error; err != nil {
			return fmt.Errorf("failed to update report %w", err)
		}
		_, err := gorm.G[models.ReportColumn](r.db).Where("report_id = ?", id).Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete report column")
		}
		for _, column := range input.Columns {
			column.ReportID = report.ID
		}
		if err := tx.WithContext(ctx).Create(&input.Columns).Error; err != nil {
			return fmt.Errorf("failed to created column report")
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to update report %w", err)
	}
	return nil
}

func (r *reportRepo) DeleteReport(ctx context.Context, id int64) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("id = ?", id).Delete(&models.Report{}).Error; err != nil {
			return fmt.Errorf("failed to delete report %w", err)
		}
		_, err := gorm.G[models.ReportColumn](r.db).Where("report_id = ?", id).Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete report column")
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to delete report %w", err)
	}
	return nil
}
func (r *reportRepo) ExportReportToExcel(ctx context.Context, reportName string, columns []models.ReportColumn, data []map[string]interface{}, fromDate, toDate time.Time) ([]byte, error) {
	// Prepare headers từ columns
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = col.Title
	}

	processedData := make([]map[string]interface{}, len(data))
	for i, row := range data {
		newRow := make(map[string]interface{})
		for _, col := range columns {
			if val, exists := row[col.Code]; exists {
				newRow[col.Title] = val
			}
		}
		processedData[i] = newRow
	}
	title := fmt.Sprintf("%s (%s to %s)",
		reportName,
		fromDate.Format("2006-01-02"),
		toDate.Format("2006-01-02"))

	filePath, err := utils.ExportToExcel(processedData, headers, title)
	if err != nil {
		return nil, fmt.Errorf("failed to export to excel: %w", err)
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read excel file: %w", err)
	}

	return fileBytes, nil
}
func (r *reportRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Report{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count reports: %w", err)
	}
	return count, nil
}
