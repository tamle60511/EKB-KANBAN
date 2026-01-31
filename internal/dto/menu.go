package dto

type MenuCreate struct {
	Title string `json:"title"`
	Code  string `json:"code"`
	Route string `json:"route,omitempty"`
	Icon  string `json:"icon,omitempty"`
}
type MenuCreateItem struct {
	Title    string `json:"title"`
	Code     string `json:"code"`
	Route    string `json:"route"`
	ReportID int64  `json:"report_id"`
}
type MenuCreateReq struct {
	DepartmentID int64            `json:"department_id"`
	Detail       MenuCreate       `json:"detail"`
	List         []MenuCreateItem `json:"list"`
}
type MenuUpdateReq struct {
	DepartmentID *int64            `json:"department_id"`
	Detail       *MenuCreate       `json:"detail"`
	List         *[]MenuUpdateItem `json:"list"`
}
type MenuUpdateItem struct {
	ID       *int64  `json:"id"`
	Title    *string `json:"title"`
	Code     *string `json:"code"`
	Route    *string `json:"route"`
	ReportID *int64  `json:"report_id"`
}

type MenuDetailReq struct {
	DepartmentID *int64 `json:"department_id" query:"department_id"`
}

type MenuDetailRes struct {
	Menu []MenuDetail `json:"menu"`
}
type MenuDetail struct {
	ID           int64            `json:"id"`
	Title        string           `json:"title"`
	Code         string           `json:"code"`
	Route        string           `json:"route"`
	Icon         string           `json:"icon"`
	DepartmentID int64            `json:"department_id"`
	List         []MenuDetailItem `json:"list"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}
type MenuDetailItem struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Code     string `json:"code"`
	Icon    string `json:"icon,omitempty"`
	Route    string `json:"route"`
	ReportID int64  `json:"report_id"`
}
