package models

type SymbolsData struct {
	BaseCurrency    string `json:"base-currency"`    // 基础币种
	QuoteCurrency   string `json:"quote-currency"`   // 计价币种
	PricePrecision  int    `json:"price-precision"`  // 价格精度位数(0为个位)
	AmountPrecision int    `json:"amount-precision"` // 数量精度位数(0为个位)
	SymbolPartition string `json:"symbol-partition"` // 交易区, main: 主区, innovation: 创新区, bifurcation: 分叉区
}

type SymbolsReturn struct {
	Status  string        `json:"status"` // 请求状态
	Data    []SymbolsData `json:"data"`   // 交易及精度数据
	ErrCode string        `json:"err-code"`
	ErrMsg  string        `json:"err-msg"`
}
