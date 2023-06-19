package common

type PaginateReq struct {
	Page     int64 `json:"page" form:"page"`           // 页数
	PageSize int64 `json:"page_size" form:"page_size"` // 每页记录数
}
