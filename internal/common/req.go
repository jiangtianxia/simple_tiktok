package common

type PaginateReq struct {
	Page     int64 `json:"page" form:"page"`           // 页数
	PageSize int64 `json:"page_size" form:"page_size"` // 每页记录数
}

type SearchReq struct {
	Where string `json:"where" form:"where"` // 条件, where=条件1:值1,条件2:值2, ......
	Order string `json:"order" form:"order"` // 排序规则, order=规则1,规则2, ......
	Sort  int    `json:"sort" form:"sort"`   // 排序方式, -1: 倒序, 1: 正序 默认倒序
}

type TokenInfoReq struct {
	Id       uint   `json:"id" form:"id"`             // 用户id
	Username string `json:"username" form:"username"` // 用户名
}
