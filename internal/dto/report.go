package dto

import (
	"cqs-kanban/internal/models"
	"time"
)

type ReportReq struct {
	ReportID int64      `json:"report_id" uri:"id"`
	FromDate *time.Time `json:"from_date,omitempty"`
	ToDate   *time.Time `json:"to_date,omitempty"`
	Period   *string    `json:"period,omitempty"`
}

type ReportRes struct {
	ReportID   int64                 `json:"report_id"`
	ReportType string                `json:"report_type"`
	ReportName string                `json:"report_name"`
	Columns    []models.ReportColumn `json:"columns"`
	Data       []map[string]any      `json:"data"`
}

type ReportCreateReq struct {
	ReportType     string          `json:"report_type"`
	ReportName     string          `json:"report_name"`
	DepartmentID   string          `json:"department_id"`
	QueryStatement string          `json:"query_statement"`
	Columns        []*ReportColumn `json:"columns"`
}

type ReportCreateModel struct {
	ReportType     string          `json:"report_type"`
	ReportName     string          `json:"report_name"`
	DepartmentID   string          `json:"department_id"`
	QueryStatement string          `json:"sql_query"`
	Columns        []*ReportColumn `json:"columns"`
}

type ReportUpdateModel struct {
	ReportType     *string         `json:"report_type"`
	ReportName     *string         `json:"report_name"`
	DepartmentID   *string         `json:"department_id"`
	QueryStatement *string         `json:"query_statement"`
	Columns        []*ReportColumn `json:"columns"`
}

type ReportColumn struct {
	ID       int64  `json:"id"`
	ReportID int64  `json:"report_id"`
	Title    string `json:"title"`
	Code     string `json:"code"`
	Type     string `json:"type"`
	Num      int64  `json:"num"`
}

type ReportFileResponse struct {
	ReportName  string    `json:"report_name"`
	FileName    string    `json:"file_name"`
	FileDetal   any       `json:"filed_detail"`
	GeneratedAt time.Time `json:"generated_at"`
}

type ReportDetail struct {
	ID             int64           `json:"id"`
	ReportType     string          `json:"report_type"`
	ReportName     string          `json:"report_name"`
	DepartmentID   string          `json:"department_id"`
	QueryStatement string          `json:"query_statement"`
	Columns        []*ReportColumn `json:"columns"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
