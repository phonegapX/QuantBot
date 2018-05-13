package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// GateIo the exchange struct of gateio.io
type GateIo struct {
	stockTypeMap     map[string]string
	tradeTypeMap     map[string]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record
	host             string
	logger           model.Logger
	option           Option

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewGateIo create an exchange struct of gateio.io
func NewGateIo(opt Option) Exchange {
	return &GateIo{
		stockTypeMap: map[string]string{
			"BTC/USDT":  "btc",
			"ETH/USDT":  "eth",
			"EOS/USDT":  "eos",
			"ONT/USDT":  "ont",
			"QTUM/USDT": "qtum",
		},
		tradeTypeMap: map[string]string{
			"buy":         constant.TradeTypeBuy,
			"sell":        constant.TradeTypeSell,
			"buy_market":  constant.TradeTypeBuy,
			"sell_market": constant.TradeTypeSell,
		},
		recordsPeriodMap: map[string]string{
			"M":   "1min",
			"M5":  "5min",
			"M15": "15min",
			"M30": "30min",
			"H":   "1hour",
			"D":   "1day",
			"W":   "1week",
		},
		minAmountMap: map[string]float64{
			"BTC/USDT":  0.001,
			"ETH/USDT":  0.001,
			"EOS/USDT":  0.001,
			"ONT/USDT":  0.001,
			"QTUM/USDT": 0.001,
		},
		records: make(map[string][]Record),
		host:    "https://data.gateio.io/api2/1/",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log print something to console
func (e *GateIo) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *GateIo) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *GateIo) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *GateIo) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *GateIo) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *GateIo) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

func (e *GateIo) getAuthJSON(url string, params []string) (json *simplejson.Json, err error) {
	e.lastTimes++
	resp, err := post_gateio(url, params, e.option.AccessKey, signSha512(params, e.option.SecretKey))
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// GetAccount get the account detail of this exchange
func (e *GateIo) GetAccount() interface{} {
	json, err := e.getAuthJSON(e.host+"private/balances", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if result := json.Get("result").MustString(); result != "true" {
		err = fmt.Errorf("the error message => %s", json.Get("message").MustString())
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	return map[string]float64{
		"USDT":       conver.Float64Must(json.GetPath("available", "USDT").Interface()),
		"FrozenUSDT": conver.Float64Must(json.GetPath("locked", "USDT").Interface()),
		"BTC":        conver.Float64Must(json.GetPath("available", "BTC").Interface()),
		"FrozenBTC":  conver.Float64Must(json.GetPath("locked", "BTC").Interface()),
		"ETH":        conver.Float64Must(json.GetPath("available", "ETH").Interface()),
		"FrozenETH":  conver.Float64Must(json.GetPath("locked", "ETH").Interface()),
		"EOS":        conver.Float64Must(json.GetPath("available", "EOS").Interface()),
		"FrozenEOS":  conver.Float64Must(json.GetPath("locked", "EOS").Interface()),
		"ONT":        conver.Float64Must(json.GetPath("available", "ONT").Interface()),
		"FrozenONT":  conver.Float64Must(json.GetPath("locked", "ONT").Interface()),
		"QTUM":       conver.Float64Must(json.GetPath("available", "QTUM").Interface()),
		"FrozenQTUM": conver.Float64Must(json.GetPath("locked", "QTUM").Interface()),
	}
}

// Trade place an order
func (e *GateIo) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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

func (e *GateIo) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	params := []string{
		"currencyPair=" + e.stockTypeMap[stockType] + "_usdt",
	}
	rateParam := fmt.Sprintf("rate=%f", price)
	amountParam := fmt.Sprintf("amount=%f", amount)
	params = append(params, rateParam, amountParam)
	json, err := e.getAuthJSON(e.host+"private/buy", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if result := json.Get("result").MustString(); result != "true" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, the error message => ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("orderNumber").Interface())
}

func (e *GateIo) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	params := []string{
		"currencyPair=" + e.stockTypeMap[stockType] + "_usdt",
	}
	rateParam := fmt.Sprintf("rate=%f", price)
	amountParam := fmt.Sprintf("amount=%f", amount)
	params = append(params, rateParam, amountParam)
	json, err := e.getAuthJSON(e.host+"private/sell", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if result := json.Get("result").MustString(); result != "true" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, the error message => ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("orderNumber").Interface())
}

// GetOrder get details of an order
func (e *GateIo) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"currencyPair=" + e.stockTypeMap[stockType] + "_usdt",
		"orderNumber=" + id,
	}
	json, err := e.getAuthJSON(e.host+"private/getOrder", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustString(); result != "true" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, the error message => ", json.Get("message").MustString())
		return false
	}
	orderJSON := json.Get("order")
	return Order{
		ID:         fmt.Sprint(orderJSON.Get("orderNumber").Interface()),
		Price:      conver.Float64Must(orderJSON.Get("rate").Interface()),
		Amount:     conver.Float64Must(orderJSON.Get("initialAmount").Interface()),
		DealAmount: conver.Float64Must(orderJSON.Get("filledAmount").Interface()),
		TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
		StockType:  stockType,
	}
}

// GetOrders get all unfilled orders
func (e *GateIo) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	json, err := e.getAuthJSON(e.host+"private/openOrders", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if result := json.Get("result").MustString(); result != "true" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, the error message => ", json.Get("message").MustString())
		return false
	}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("orderNumber").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("initialRate").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("initialAmount").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("filledAmount").Interface()),
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *GateIo) GetTrades(stockType string) interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *GateIo) CancelOrder(order Order) bool {
	params := []string{
		"currencyPair=" + e.stockTypeMap[order.StockType] + "_usdt",
		"orderNumber=" + order.ID,
	}
	json, err := e.getAuthJSON(e.host+"private/cancelOrder", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, the error message => ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *GateIo) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	resp, err := get(fmt.Sprintf("http://data.gateio.io/api2/1/orderBook/%v_usdt", e.stockTypeMap[stockType]))
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	depthsJSON := json.Get("bids")
	for i := 0; i < len(depthsJSON.MustArray()); i++ {
		depthJSON := depthsJSON.GetIndex(i)
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	depthsJSON = json.Get("asks")
	for i := len(depthsJSON.MustArray()); i > 0; i-- {
		depthJSON := depthsJSON.GetIndex(i - 1)
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
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
func (e *GateIo) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *GateIo) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return nil
}
