package repository

import (
	"context"
	"cqs-kanban/internal/dto"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type forecastRepo struct {
	db *gorm.DB
}

type ForecastRepo interface {
	GetForecast(ctx context.Context, req *dto.ReportReq) ([]dto.CombinedForecast, error)
	ExportGroupedForecastToExcel(ctx context.Context, data []dto.CombinedForecast, fromDate, toDate time.Time) ([]byte, error)
}

func NewForecastRepo(db *gorm.DB) ForecastRepo {
	return &forecastRepo{db: db}
}

func (r *forecastRepo) GetForecast(ctx context.Context, req *dto.ReportReq) ([]dto.CombinedForecast, error) {
	query := `WITH TD AS (
    SELECT
        RTRIM(TD001) + '-' + RTRIM(TD002) + '-' + RTRIM(TD003) AS TD01,
		TC004 AS MKH,
        MA002 AS KH01,
        TD004 AS TD02,
        TD005 AS TD03,
        TD006 AS TD04,
        TD008 AS TD05,
        TD009 AS TD06,
        TD011 AS TD07,
        CONVERT(VARCHAR(10), CONVERT(DATETIME, TD013), 103) AS TD08
    FROM
        COPTD
    LEFT JOIN COPTC ON TC001 = TD001 AND TC002 = TD002
    LEFT JOIN COPMA ON MA001 = TC004
    WHERE TC039 >= ? AND TC039 <= ?
),
MF AS (
    SELECT
        RTRIM(MF001) AS TD09,
        MA002 AS KH02,
        MF003 AS TD10,
        MF004 AS TD11,
        MF005 AS TD12,
        MF008 AS TD13,
        MF012 AS TD14,
        CASE
            WHEN ISDATE(MF006) = 1 THEN CONVERT(VARCHAR(10), CONVERT(DATETIME, MF006), 103)
            ELSE NULL
        END AS TD15
    FROM
        COPMF
    LEFT JOIN COPME ON ME001 = MF001
    LEFT JOIN COPMA ON MA001 = ME002
)


SELECT
    TD.TD01,
	TD.MKH,
    TD.KH01,
    TD.TD02,
    TD.TD03,
    TD.TD04,
    TD.TD05,
    TD.TD06,
    TD.TD07,
    TD.TD08,
    NULL AS TD09,
	NULL AS KH02,
    NULL AS TD10,
    NULL AS TD11,
    NULL AS TD12,
    NULL AS TD13,
    NULL AS TD14,
    NULL AS TD15
FROM
    TD

UNION ALL

SELECT
    NULL,
    NULL,
	 NULL,
    NULL,
    NULL,
    NULL, 
    NULL, 
    NULL, 
    NULL, 
    NULL,
    MF.TD09,
	MF.KH02,
    MF.TD10,
    MF.TD11,
    MF.TD12,
    MF.TD13,
    MF.TD14,
    MF.TD15
FROM
    MF
ORDER BY
    TD01, TD02;`

	var forecasts []dto.Forecast
	result := r.db.WithContext(ctx).Raw(query, req.FromDate, req.ToDate).Scan(&forecasts)
	if result.Error != nil {
		return nil, result.Error
	}

	groupedForecasts := make(map[string]*dto.CombinedForecast)
	groupedDetails := make(map[string][]dto.SubForecast)

	for _, forecast := range forecasts {
		td02 := forecast.TD02
		td10 := forecast.TD10

		if forecast.TD01 != "" {
			if group, exists := groupedForecasts[td02]; exists {
				group.Columns = append(group.Columns, dto.Columns{
					TD01: forecast.TD01, MKH: forecast.MKH, KH01: forecast.KH01, TD02: td02, TD03: forecast.TD03, TD04: forecast.TD04,
					TD05: forecast.TD05, TD06: forecast.TD06, TD07: forecast.TD07, TD08: forecast.TD08,
				})
			} else {
				groupedForecasts[td02] = &dto.CombinedForecast{
					Columns: []dto.Columns{{
						TD01: forecast.TD01, MKH: forecast.MKH, KH01: forecast.KH01, TD02: td02, TD03: forecast.TD03, TD04: forecast.TD04,
						TD05: forecast.TD05, TD06: forecast.TD06, TD07: forecast.TD07, TD08: forecast.TD08,
					}},
					Details: []dto.SubForecast{},
				}
			}
		} else if forecast.TD09 != "" {
			detail := dto.SubForecast{
				TD09: forecast.TD09, KH02: forecast.KH02, TD10: td10, TD11: forecast.TD11, TD12: forecast.TD12,
				TD13: forecast.TD13, TD14: forecast.TD14, TD15: forecast.TD15,
			}
			groupedDetails[td10] = append(groupedDetails[td10], detail)
		}
	}

	for td02, group := range groupedForecasts {
		if details, exists := groupedDetails[td02]; exists {
			group.Details = details
		}
	}

	combinedForecasts := make([]dto.CombinedForecast, 0, len(groupedForecasts))
	for _, group := range groupedForecasts {
		if len(group.Details) > 0 {
			combinedForecasts = append(combinedForecasts, *group)
		}
	}

	return combinedForecasts, nil
}

func (r *forecastRepo) ExportGroupedForecastToExcel(ctx context.Context, data []dto.CombinedForecast, fromDate, toDate time.Time) ([]byte, error) {  
	f := excelize.NewFile()  
	defer f.Close()  

	sheetName := "Forecast Report"  
	index, err := f.NewSheet(sheetName)  
	if err != nil {  
		return nil, fmt.Errorf("failed to create sheet: %w", err)  
	}  
	f.SetActiveSheet(index)  
	f.DeleteSheet("Sheet1")  

	// ===== STYLES =====  
	titleStyle, _ := f.NewStyle(&excelize.Style{  
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "FFFFFF"},  
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1F4E78"}, Pattern: 1},  
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},  
		Border: []excelize.Border{  
			{Type: "top", Color: "000000", Style: 2},  
			{Type: "bottom", Color: "000000", Style: 2},  
			{Type: "left", Color: "000000", Style: 2},  
			{Type: "right", Color: "000000", Style: 2},  
		},  
	})  

	groupHeaderStyle, _ := f.NewStyle(&excelize.Style{  
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},  
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},  
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},  
		Border: []excelize.Border{  
			{Type: "top", Color: "000000", Style: 2},  
			{Type: "bottom", Color: "000000", Style: 2},  
			{Type: "left", Color: "000000", Style: 2},  
			{Type: "right", Color: "000000", Style: 2},  
		},  
	})  

	sectionHeaderStyle, _ := f.NewStyle(&excelize.Style{  
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "000000"},  
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"D9E1F2"}, Pattern: 1},  
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},  
		Border: []excelize.Border{  
			{Type: "top", Color: "000000", Style: 1},  
			{Type: "bottom", Color: "000000", Style: 1},  
			{Type: "left", Color: "000000", Style: 1},  
			{Type: "right", Color: "000000", Style: 1},  
		},  
	})  

	orderHeaderStyle, _ := f.NewStyle(&excelize.Style{  
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "000000"},  
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"E2EFDA"}, Pattern: 1},  
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},  
		Border: []excelize.Border{  
			{Type: "top", Color: "000000", Style: 1},  
			{Type: "bottom", Color: "000000", Style: 1},  
			{Type: "left", Color: "000000", Style: 1},  
			{Type: "right", Color: "000000", Style: 1},  
		},  
	})  

	cellStyle, _ := f.NewStyle(&excelize.Style{  
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},  
	})  

	numberStyle, _ := f.NewStyle(&excelize.Style{  
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},  
		NumFmt:    2, // 0.00  
	})  

	dateStyle, _ := f.NewStyle(&excelize.Style{  
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},  
	})  
	  
	maxCols := "J" // The widest column used is J for Order Details  

	// ===== REPORT TITLE =====  
	currentRow := 1  
	titleText := fmt.Sprintf("FORECAST REPORT (%s to %s)", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))  
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), titleText)  
	f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("%s%d", maxCols, currentRow))  
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("%s%d", maxCols, currentRow), titleStyle)  
	f.SetRowHeight(sheetName, currentRow, 30)  
	currentRow += 2  

	// ===== PROCESS EACH GROUP =====  
	for groupIndex, group := range data {  
		if len(group.Columns) == 0 {  
			continue  
		}  

		orderNumber := group.Columns[0].TD02  

		// --- GROUP HEADER ---  
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("Group %d - Order: %s", groupIndex+1, orderNumber))  
		f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("%s%d", maxCols, currentRow))  
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("%s%d", maxCols, currentRow), groupHeaderStyle)  
		f.SetRowHeight(sheetName, currentRow, 25)  
		currentRow++  

		// --- ORDER DETAILS SECTION ---  
		if len(group.Columns) > 0 {  
			// Section Label  
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), "ORDER DETAILS")  
			f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("%s%d", maxCols, currentRow))  
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("%s%d", maxCols, currentRow), sectionHeaderStyle)  
			currentRow++  

			// Headers  
			orderHeaders := []string{"Đơn hàng", "Mã KH", "Tên KH", "Mã Sản Phẩm", "Tên Sản Phẩm", "Quy cách", "Số lượng đặt", "Số lượng đã giao", "Đơn giá", "Ngày dự định giao"}  
			for i, header := range orderHeaders {  
				cell := fmt.Sprintf("%s%d", string(rune('A'+i)), currentRow)  
				f.SetCellValue(sheetName, cell, header)  
				f.SetCellStyle(sheetName, cell, cell, sectionHeaderStyle)  
			}  
			currentRow++  

			// Data Rows  
			for _, col := range group.Columns {  
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), col.TD01)  
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), col.MKH)  
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), col.KH01)  
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), col.TD02)  
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), col.TD03)  
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), col.TD04)  
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", currentRow), col.TD05)  
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", currentRow), col.TD06)  
				f.SetCellValue(sheetName, fmt.Sprintf("I%d", currentRow), col.TD07)  
				f.SetCellValue(sheetName, fmt.Sprintf("J%d", currentRow), col.TD08)  

				// Apply styles  
				textCols := []string{"A", "B", "C", "D", "E", "F"}  
				for _, c := range textCols {  
					f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, currentRow), fmt.Sprintf("%s%d", c, currentRow), cellStyle)  
				}  
				numCols := []string{"G", "H", "I"}  
				for _, c := range numCols {  
					f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, currentRow), fmt.Sprintf("%s%d", c, currentRow), numberStyle)  
				}  
				f.SetCellStyle(sheetName, fmt.Sprintf("J%d", currentRow), fmt.Sprintf("J%d", currentRow), dateStyle)  

				currentRow++  
			}  
			currentRow++ // Blank row  
		}  

		// --- FORECAST SCHEDULE SECTION ---  
		if len(group.Details) > 0 {  
			// Section Label  
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("FORECAST SCHEDULE (%d records)", len(group.Details)))  
			f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("H%d", currentRow)) // Merge up to H  
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("H%d", currentRow), orderHeaderStyle)  
			currentRow++  

			// Headers  
			forecastHeaders := []string{"Mã dự đoán", "Tên KH", "Mã Sản Phẩm", "Tên Sản Phẩm", "Quy cách", "Số lượng đặt", "Đơn giá", "Ngày"}  
			for i, header := range forecastHeaders {  
				cell := fmt.Sprintf("%s%d", string(rune('A'+i)), currentRow)  
				f.SetCellValue(sheetName, cell, header)  
				f.SetCellStyle(sheetName, cell, cell, orderHeaderStyle)  
			}  
			currentRow++  

			// Data Rows  
			for _, detail := range group.Details {  
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), detail.TD09)  
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), detail.KH02)  
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), detail.TD10)  
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), detail.TD11)  
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), detail.TD12)  
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), detail.TD13)  
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", currentRow), detail.TD14)  
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", currentRow), detail.TD15)  

				// Apply styles  
				textCols := []string{"A", "B", "C", "D", "E"}  
				for _, c := range textCols {  
					f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, currentRow), fmt.Sprintf("%s%d", c, currentRow), cellStyle)  
				}  
				numCols := []string{"F", "G"}  
				for _, c := range numCols {  
					f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, currentRow), fmt.Sprintf("%s%d", c, currentRow), numberStyle)  
				}  
				f.SetCellStyle(sheetName, fmt.Sprintf("H%d", currentRow), fmt.Sprintf("H%d", currentRow), dateStyle)  

				currentRow++  
			}  
			currentRow += 2  
		}  
	}  

	columnWidths := map[string]float64{  
		"A": 20, "B": 20, "C": 30, "D": 20, "E": 30, "F": 25, "G": 15, "H": 18, "I": 15, "J": 18,  
	}  
	for col, width := range columnWidths {  
		f.SetColWidth(sheetName, col, col, width)  
	}  

	buffer, err := f.WriteToBuffer()  
	if err != nil {  
		return nil, fmt.Errorf("failed to write to buffer: %w", err)  
	}  

	return buffer.Bytes(), nil  
}  