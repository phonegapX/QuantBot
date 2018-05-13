package models

type TradeData struct {
	ID        int64   `json:"id"`        //成交ID
	Price     float64 `json:"price"`     // 成交价
	Amount    float64 `json:"amount"`    // 成交量
	Direction string  `json:"direction"` // 主动成交方向
	Ts        int64   `json:"ts"`        // 成交时间
}

type TradeTick struct {
	ID   int64       `json:"id"`   // 消息ID
	Ts   int64       `json:"ts"`   // 最新成交时间
	Data []TradeData `json:"data"` // Trade数据
}

type TradeReturn struct {
	Status  string      `json:"status"` // 请求状态, ok或者error
	Ch      string      `json:"ch"`     // 数据所属的Channel, 格式: market.$symbol.trade.detail
	Ts      int64       `json:"ts"`     // 发送时间
	Data    []TradeTick `json:"data"`   // 成交记录
	ErrCode string      `json:"err-code"`
	ErrMsg  string      `json:"err-msg"`
}
