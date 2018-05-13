package models

type KLineData struct {
	ID     int64   `json:"id"`     // K线ID
	Amount float64 `json:"amount"` // 成交量
	Count  int64   `json:"count"`  // 成交笔数
	Open   float64 `json:"open"`   // 开盘价
	Close  float64 `json:"close"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"low"`    // 最低价
	High   float64 `json:"high"`   // 最高价
	Vol    float64 `json:"vol"`    // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
}

type KLineReturn struct {
	Status  string      `json:"status"`   // 请求处理结果, "ok"、"error"
	Ts      int64       `json:"ts"`       // 响应生成时间点, 单位毫秒
	Data    []KLineData `json:"data"`     // KLine数据
	Ch      string      `json:"ch"`       // 数据所属的Channel, 格式: market.$symbol.kline.$period
	ErrCode string      `json:"err-code"` // 错误代码
	ErrMsg  string      `json:"err-msg"`  // 错误提示
}
