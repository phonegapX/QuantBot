package models

type OrderDetail struct {
	ID         int64  `json:"id"`           //订单ID
	Price      string `json:"price"`        //价格
	Amount     string `json:"amount"`       //总量
	DealAmount string `json:"field-amount"` //成交量
	TradeType  string `json:"type"`         //交易类型
	StockType  string `json:"symbol"`       //货币类型
}

type OrderDetailReturn struct {
	Status  string      `json:"status"` // 请求状态
	Data    OrderDetail `json:"data"`   // 订单详情
	ErrCode string      `json:"err-code"`
	ErrMsg  string      `json:"err-msg"`
}
