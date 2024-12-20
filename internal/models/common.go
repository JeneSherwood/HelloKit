package models

type Paging struct {
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_size"`
	Bookmark string `json:"bookmark"` // CouchDB书签
}

type ListReq struct {
	Paging    `json:"paging"`
	Sort      string `json:"sort"`
	SortField string `json:"sort_field"`
	Account   string `json:"account"`   // 客户账号
	FundCode  string `json:"fund_code"` // 产品代码
	TtmID     string `json:"ttm_id"`    // 事务id
	State     int    `json:"state"`     // 状态
}

type Request struct {
	Payload []byte `json:"payload"` //请求业务报文内容体，由业务系统定义
}
