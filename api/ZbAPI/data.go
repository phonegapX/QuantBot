package ZbAPI

//==================================================//
type depthOrder []float64
type respDepth struct {
	Timestamp int          `mapstructure:"timestamp"`
	Asks      []depthOrder `mapstructure:"asks"`
	Bids      []depthOrder `mapstructure:"bids"`
}

//==================================================//
type klineData []float64
type respKline struct {
	Symbol    string      `json:"symbol"`
	Data      []klineData `json:"data"`
	MoneyType string      `json:"moneyType"`
}

//==================================================//
type respTrades []trade
type trade struct {
	Amount    string `json:"amount"`
	Date      int64  `json:"date"`
	Price     string `json:"price"`
	Tid       int64  `json:"tid"`
	TradeType string `json:"trade_type"`
	TxType    string `json:"type"`
}

//==================================================//
type ticker struct {
	Vol  string `json:"vol"`
	Last string `json:"last"`
	Sell string `json:"sell"`
	Buy  string `json:"buy"`
	High string `json:"high"`
	Low  string `json:"low"`
}

type respTicker struct {
	Ticker ticker `json:"ticker"`
	Date   string `json:"date"`
}

//==================================================//
type accountInfoCoin struct {
	EnName        string `mapstructure:"enName"`
	Freez         string `mapstructure:"freez"`
	UnitDecimal   int    `mapstructure:"unitDecimal"`
	CnName        string `mapstructure:"cnName"`
	IsCanRecharge bool   `mapstructure:"isCanRecharge"`
	UnitTag       string `mapstructure:"unitTag"`
	IsCanWithdraw bool   `mapstructure:"isCanWithdraw"`
	Available     string `mapstructure:"available"`
	Key           string `mapstructure:"key"`
}
type accountInfoBase struct {
	Username             string `mapstructure:"username"`
	TradePasswordEnabled bool   `mapstructure:"trade_password_enabled"`
	AuthGoogleEnabled    bool   `mapstructure:"auth_google_enabled"`
	AuthMobileEnabled    bool   `mapstructure:"auth_mobile_enabled"`
}
type accountInfoResult struct {
	Coins    []accountInfoCoin `mapstructure:"coins"`
	BaseInfo accountInfoBase   `mapstructure:"base"`
}
type respAccountInfo struct {
	Result  accountInfoResult `mapstructure:"result"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
}

//==================================================//
type respOrders []order
type order struct {
	Currency    string  `json:"currency"`
	ID          string  `json:"id"`
	Price       float64 `json:"price"`
	Status      int     `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	TradeAmount float64 `json:"trade_amount"`
	TradeDate   int64   `json:"trade_date"`
	TradeMoney  string  `json:"trade_money"`
	TradePrice  float64 `json:"trade_price"`
	OrderType   int     `json:"type"`
	Code        int     `json:"code"`
	Message     string  `json:"message"`
}

//==================================================//
type respSimple struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//==================================================//
type respOrder struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

//==================================================//
