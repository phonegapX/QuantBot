package models

type MarketDepth struct {
	ID   int64       `json:"id"`   // 消息ID
	Ts   int64       `json:"ts"`   // 消息声称事件, 单位: 毫秒
	Bids [][]float64 `json:"bids"` // 买盘, [price(成交价), amount(成交量)], 按price降序排列
	Asks [][]float64 `json:"asks"` // 卖盘, [price(成交价), amount(成交量)], 按price升序排列
}

type MarketDepthReturn struct {
	Status  string      `json:"status"` // 请求状态, ok或者error
	Ts      int64       `json:"ts"`     // 响应生成时间点, 单位: 毫秒
	Tick    MarketDepth `json:"tick"`   // Depth数据
	Ch      string      `json:"ch"`     //  数据所属的Channel, 格式: market.$symbol.depth.$type
	ErrCode string      `json:"err-code"`
	ErrMsg  string      `json:"err-msg"`
}
