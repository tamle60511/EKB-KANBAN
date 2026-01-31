package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportToExcel exports data to Excel file
func ExportToExcel(data []map[string]interface{}, headers []string, title string) (string, error) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Get the default sheet
	sheetName := "Sheet1"

	// Set title
	f.SetCellValue(sheetName, "A1", title)

	// Set title style
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  16,
			Bold:  true,
			Color: "1F497D",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating title style: %w", err)
	}

	// Apply title style and merge cells for title
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", rune('A'+len(headers)-1)), titleStyle)
	f.MergeCell(sheetName, "A1", fmt.Sprintf("%c1", rune('A'+len(headers)-1)))

	// Set headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating header style: %w", err)
	}

	// Write headers
	for i, header := range headers {
		cellPos := fmt.Sprintf("%c3", rune('A'+i))
		f.SetCellValue(sheetName, cellPos, header)
	}

	// Apply header style
	headerRange := fmt.Sprintf("A3:%c3", rune('A'+len(headers)-1))
	f.SetCellStyle(sheetName, headerRange, headerRange, headerStyle)

	// Data cell styles
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating data style: %w", err)
	}

	// Number format style
	numberStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		NumFmt: 3, // #,##0 format
	})
	if err != nil {
		return "", fmt.Errorf("error creating number style: %w", err)
	}

	// Write data
	for i, item := range data {
		row := i + 4 // Data starts from row 4

		for j, header := range headers {
			cellPos := fmt.Sprintf("%c%d", rune('A'+j), row)
			f.SetCellValue(sheetName, cellPos, item[header])

			// Apply style based on data type
			switch header {
			case "Mã sản phẩm", "Tên sản phẩm", "Đơn vị":
				f.SetCellStyle(sheetName, cellPos, cellPos, dataStyle)
			default:
				f.SetCellStyle(sheetName, cellPos, cellPos, numberStyle)
			}
		}
	}

	// Set column width
	for i := range headers {
		colName := string(rune('A' + i))
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Set row height
	f.SetRowHeight(sheetName, 1, 30)
	f.SetRowHeight(sheetName, 3, 25)

	// Generate timestamp for filename
	timestamp := time.Now().Format("20060102_150405")

	// Create downloads directory if it doesn't exist
	downloadPath := "public/downloads"
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return "", fmt.Errorf("error creating downloads directory: %w", err)
	}

	// Create sanitized filename
	safeTitlePart := sanitizeFilename(title)
	if len(safeTitlePart) > 30 {
		safeTitlePart = safeTitlePart[:30]
	}

	// Complete filename
	filename := fmt.Sprintf("%s_%s.xlsx", safeTitlePart, timestamp)
	filePath := filepath.Join(downloadPath, filename)

	// Save the file
	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("error saving Excel file: %w", err)
	}

	return filePath, nil
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(name string) string {
	name = filepath.Clean(name)
	invalidChars := []rune{'<', '>', ':', '"', '/', '\\', '|', '?', '*'}

	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, string(char), "_")
	}

	return name
}
