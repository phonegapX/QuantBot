package models

type PlaceRequestParams struct {
	AccountID string `json:"account-id"` // 账户ID
	Amount    string `json:"amount"`     // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
	Price     string `json:"price"`      // 下单价格, 市价单不传该参数
	Source    string `json:"source"`     // 订单来源, api: API调用, margin-api: 借贷资产交易
	Symbol    string `json:"symbol"`     // 交易对, btcusdt, bccbtc......
	Type      string `json:"type"`       // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
}

type PlaceReturn struct {
	Status  string `json:"status"`
	Data    string `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}
