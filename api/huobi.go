package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/api/HuobiProAPI/config"
	"github.com/phonegapX/QuantBot/api/HuobiProAPI/models"
	"github.com/phonegapX/QuantBot/api/HuobiProAPI/services"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// Huobi the exchange struct of huobi.com
type Huobi struct {
	stockTypeMap     map[string]string
	tradeTypeMap     map[string]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record
	logger           model.Logger
	option           Option

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewHuobi create an exchange struct of huobi.com
func NewHuobi(opt Option) Exchange {
	config.ACCESS_KEY = opt.AccessKey
	config.SECRET_KEY = opt.SecretKey
	//...
	return &Huobi{
		stockTypeMap: map[string]string{
			"BTC/USDT":  "btc",
			"ETH/USDT":  "eth",
			"EOS/USDT":  "eos",
			"ONT/USDT":  "ont",
			"QTUM/USDT": "qtum",
		},
		tradeTypeMap: map[string]string{
			"buy-limit":   constant.TradeTypeBuy,
			"sell-limit":  constant.TradeTypeSell,
			"buy-market":  constant.TradeTypeBuy,
			"sell-market": constant.TradeTypeSell,
		},
		recordsPeriodMap: map[string]string{
			"M":   "001",
			"M5":  "005",
			"M15": "015",
			"M30": "030",
			"H":   "060",
			"D":   "100",
			"W":   "200",
		},
		minAmountMap: map[string]float64{
			"BTC/USDT":  0.001,
			"ETH/USDT":  0.001,
			"EOS/USDT":  0.001,
			"ONT/USDT":  0.001,
			"QTUM/USDT": 0.001,
		},
		records: make(map[string][]Record),
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log print something to console
func (e *Huobi) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *Huobi) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *Huobi) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *Huobi) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *Huobi) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *Huobi) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *Huobi) GetAccount() interface{} {
	accounts, err := services.GetAccounts()
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if accounts.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", accounts.ErrMsg)
		return false
	}
	accountID := int64(-1)
	count := len(accounts.Data)
	for i := 0; i < count; i++ {
		actData := accounts.Data[i]
		if actData.State == "working" && actData.Type == "spot" {
			accountID = actData.ID
			break
		}
	}
	if accountID == -1 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", "all account locked")
		return false
	}
	balance, err := services.GetAccountBalance(strconv.FormatInt(accountID, 10))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if balance.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", balance.ErrMsg)
		return false
	}
	result := make(map[string]float64)
	count = len(balance.Data.List)
	for i := 0; i < count; i++ {
		subAcc := balance.Data.List[i]
		if subAcc.Type == "trade" {
			result[strings.ToUpper(subAcc.Currency)] = conver.Float64Must(subAcc.Balance)
		} else if subAcc.Type == "frozen" {
			result["Frozen"+strings.ToUpper(subAcc.Currency)] = conver.Float64Must(subAcc.Balance)
		}
	}
	//...
	config.ACCOUNT_ID = strconv.FormatInt(accountID, 10)
	//...
	return result
}

// Trade place an order
func (e *Huobi) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	tradeType = strings.ToUpper(tradeType)
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized stockType: ", stockType)
		return false
	}
	switch tradeType {
	case constant.TradeTypeBuy:
		return e.buy(stockType, price, amount, msgs...)
	case constant.TradeTypeSell:
		return e.sell(stockType, price, amount, msgs...)
	default:
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized tradeType: ", tradeType)
		return false
	}
}

func (e *Huobi) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	params := models.PlaceRequestParams{
		AccountID: config.ACCOUNT_ID,                  // 账户ID
		Amount:    conver.StringMust(amount),          // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
		Price:     conver.StringMust(price),           // 下单价格, 市价单不传该参数
		Source:    "api",                              // 订单来源, api: API调用, margin-api: 借贷资产交易
		Symbol:    e.stockTypeMap[stockType] + "usdt", // 交易对, btcusdt, bccbtc......
		Type:      "buy-limit",                        // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
	}
	result, err := services.Place(params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if result.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", result.ErrMsg)
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return result.Data
}

func (e *Huobi) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	params := models.PlaceRequestParams{
		AccountID: config.ACCOUNT_ID,                  // 账户ID
		Amount:    conver.StringMust(amount),          // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
		Price:     conver.StringMust(price),           // 下单价格, 市价单不传该参数
		Source:    "api",                              // 订单来源, api: API调用, margin-api: 借贷资产交易
		Symbol:    e.stockTypeMap[stockType] + "usdt", // 交易对, btcusdt, bccbtc......
		Type:      "sell-limit",                       // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
	}
	result, err := services.Place(params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if result.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", result.ErrMsg)
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return result.Data
}

// GetOrder get details of an order
func (e *Huobi) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := services.GetOrderDetail(id)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if result.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", result.ErrMsg)
		return false
	}
	return Order{
		ID:         fmt.Sprint(result.Data.ID),
		Price:      conver.Float64Must(result.Data.Price),
		Amount:     conver.Float64Must(result.Data.Amount),
		DealAmount: conver.Float64Must(result.Data.DealAmount),
		TradeType:  e.tradeTypeMap[result.Data.TradeType],
		StockType:  stockType,
	}
}

// GetOrders get all unfilled orders
func (e *Huobi) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := services.GetOrders(e.stockTypeMap[stockType] + "usdt")
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if result.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", result.ErrMsg)
		return false
	}
	orders := []Order{}
	count := len(result.Data)
	for i := 0; i < count; i++ {
		orders = append(orders, Order{
			ID:         fmt.Sprint(result.Data[i].ID),
			Price:      conver.Float64Must(result.Data[i].Price),
			Amount:     conver.Float64Must(result.Data[i].Amount),
			DealAmount: conver.Float64Must(result.Data[i].DealAmount),
			TradeType:  e.tradeTypeMap[result.Data[i].TradeType],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *Huobi) GetTrades(stockType string) interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *Huobi) CancelOrder(order Order) bool {
	result, err := services.SubmitCancel(order.ID)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result.Status != "ok" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", result.ErrMsg)
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *Huobi) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	result, err := services.GetMarketDepth(e.stockTypeMap[stockType]+"usdt", "step0")
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	if result.Status != "ok" {
		err = fmt.Errorf("GetTicker() error, %+v", result.ErrMsg)
		return
	}
	count := len(result.Tick.Bids)
	for i := 0; i < count; i++ {
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  result.Tick.Bids[i][0],
			Amount: result.Tick.Bids[i][1],
		})
	}
	count = len(result.Tick.Asks)
	for i := 0; i < count; i++ {
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  result.Tick.Asks[i][0],
			Amount: result.Tick.Asks[i][1],
		})
	}
	if len(ticker.Bids) < 1 || len(ticker.Asks) < 1 {
		err = fmt.Errorf("GetTicker() error, can not get enough Bids or Asks")
		return
	}
	ticker.Buy = ticker.Bids[0].Price
	ticker.Sell = ticker.Asks[0].Price
	ticker.Mid = (ticker.Buy + ticker.Sell) / 2
	return
}

// GetTicker get market ticker & depth
func (e *Huobi) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *Huobi) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return nil
}
