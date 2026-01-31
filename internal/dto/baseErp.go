package dto

import "time"

type BaseERPReq struct {
	SqlQuery string
	FromDate time.Time
	ToDate   time.Time
}

type BaseERP struct {
	Data []map[string]any
}
