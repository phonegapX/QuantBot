package models

type TradeDetailData struct {
	ID        int64   `json:"id"`        // 成交ID
	Price     float64 `json:"price"`     // 成交价
	Amount    float64 `json:"amount"`    // 成交量
	Direction string  `json:"direction"` // 主动成交方向
	Ts        int64   `json:"ts"`        // 成交时间
}

type TradeDetail struct {
	ID   int64             `json:"id"`   // 消息ID
	Ts   int64             `json:"ts"`   // 最新成交时间
	Data []TradeDetailData `json:"data"` // 交易细节数据
}

type TradeDetailReturn struct {
	Status  string      `json:"status"`   // 请求处理结果, "ok"、"error"
	Ts      int64       `json:"ts"`       // 响应生成时间点, 单位毫秒
	Tick    TradeDetail `json:"tick"`     // TradeDetail数据
	Ch      string      `json:"ch"`       // 数据所属的Channel, 格式: market.$symbol.trade.detail
	ErrCode string      `json:"err-code"` // 错误代码
	ErrMsg  string      `json:"err-msg"`  // 错误提示
}
