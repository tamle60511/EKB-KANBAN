package dto

type DepartmenCreateReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type DepartmentUpdateReq struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
	Desc *string `json:"desc"`
}

type DepartmentRes struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Desc      string `json:"desc"`
	CreatedAt string `json:"created_at"`
}
