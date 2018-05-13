package models

type OrdersReturn struct {
	Status  string        `json:"status"` // 请求状态
	Data    []OrderDetail `json:"data"`   // 订单列表
	ErrCode string        `json:"err-code"`
	ErrMsg  string        `json:"err-msg"`
}
