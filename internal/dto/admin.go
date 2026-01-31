package dto

import "time"

// Operation DTOs
type OperationCreateReq struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=500"`
}

type OperationUpdateReq struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// Analytics DTOs
type DashboardStatsResponse struct {
	TotalOperations  int64                `json:"total_operations"`
	TotalDepartments int64                `json:"total_departments"`
	TotalUsers       int64                `json:"total_users"`
	TotalReports     int64                `json:"total_reports"`
	TodayAccess      int64                `json:"today_access"`
	FailedAttempts   int64                `json:"failed_attempts"`
	TopOperations    []*TopOperationStat  `json:"top_operations"`
	TopUsers         []*TopUserStat       `json:"top_users"`
	RecentLogs       []*AccessLogResponse `json:"recent_logs"`
	AccessTrend      []*AccessTrendData   `json:"access_trend"`
}

type TopOperationStat struct {
	OperationID   int    `json:"operation_id"`
	OperationName string `json:"operation_name"`
	OperationCode string `json:"operation_code"`
	AccessCount   int64  `json:"access_count"`
}

type TopUserStat struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	AccessCount int64  `json:"access_count"`
}

type AccessTrendData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type AccessLogQueryReq struct {
	UserID      int64     `query:"user_id"`
	OperationID int       `query:"operation_id"`
	Status      string    `query:"status"`
	FromDate    time.Time `query:"from_date"`
	ToDate      time.Time `query:"to_date"`
	Page        int       `query:"page" validate:"min=1"`
	PageSize    int       `query:"page_size" validate:"min=1,max=100"`
}

type AccessLogListResponse struct {
	Logs       []*AccessLogResponse `json:"logs"`
	TotalCount int64                `json:"total_count"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// User Activity Report
type UserActivityReportReq struct {
	UserID   int64     `query:"user_id" validate:"required"`
	FromDate time.Time `query:"from_date" validate:"required"`
	ToDate   time.Time `query:"to_date" validate:"required"`
}

type UserActivityReportResponse struct {
	UserID             int64                     `json:"user_id"`
	Username           string                    `json:"username"`
	TotalAccess        int64                     `json:"total_access"`
	SuccessCount       int64                     `json:"success_count"`
	FailedCount        int64                     `json:"failed_count"`
	UniqueOperations   int                       `json:"unique_operations"`
	MostUsedOperation  *TopOperationStat         `json:"most_used_operation"`
	DailyActivity      []*DailyActivityStat      `json:"daily_activity"`
	OperationBreakdown []*OperationBreakdownStat `json:"operation_breakdown"`
}

type DailyActivityStat struct {
	Date         string `json:"date"`
	AccessCount  int64  `json:"access_count"`
	SuccessCount int64  `json:"success_count"`
	FailedCount  int64  `json:"failed_count"`
}

type OperationBreakdownStat struct {
	OperationID   int    `json:"operation_id"`
	OperationName string `json:"operation_name"`
	AccessCount   int64  `json:"access_count"`
}

// Security Alert DTO
type SecurityAlertResponse struct {
	AlertType    string    `json:"alert_type"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	IPAddress    string    `json:"ip_address"`
	FailedCount  int64     `json:"failed_count"`
	LastFailedAt time.Time `json:"last_failed_at"`
	Description  string    `json:"description"`
}
