package models

type MarketDetail struct {
	ID     int64   `json:"id"`     // 消息ID
	Ts     int64   `json:"ts"`     // 24小时统计时间
	Amount float64 `json:"amount"` // 24小时成交量
	Open   float64 `json:"open"`   // 前24小时成交价
	Close  float64 `json:"close"`  // 当前成交价
	High   float64 `json:"high"`   // 近24小时最高价
	Low    float64 `json:"low"`    // 近24小时最低价
	Count  int64   `json:"count"`  // 近24小时累计成交数
	Vol    float64 `json:"vol"`    // 近24小时累计成交额, 即SUM(每一笔成交价 * 该笔的成交量)
}

type MarketDetailReturn struct {
	Status  string       `json:"status"` // 请求状态
	Ts      int64        `json:"ts"`     // 响应生成时间点
	Tick    MarketDetail `json:"tick"`   // Market Detail 24小时成交量数据
	Ch      string       `json:"ch"`     // 数据所属的Channel, 格式: market.$symbol.depth.$type
	ErrCode string       `json:"err-code"`
	ErrMsg  string       `json:"err-msg"`
}
