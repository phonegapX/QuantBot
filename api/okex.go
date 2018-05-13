package api

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// OKEX the exchange struct of okex.com
type OKEX struct {
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

// NewOKEX create an exchange struct of okex.com
func NewOKEX(opt Option) Exchange {
	return &OKEX{
		stockTypeMap: map[string]string{
			"BTC/USDT":  "btc_usdt",
			"ETH/USDT":  "eth_usdt",
			"EOS/USDT":  "eos_usdt",
			"ONT/USDT":  "ont_usdt",
			"QTUM/USDT": "qtum_usdt",
			"ONT/ETH":   "ont_eth",
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
			"ONT/ETH":   0.001,
		},
		records: make(map[string][]Record),
		host:    "https://www.okex.com/api/v1/",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log print something to console
func (e *OKEX) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *OKEX) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *OKEX) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *OKEX) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *OKEX) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *OKEX) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

func (e *OKEX) getAuthJSON(url string, params []string) (json *simplejson.Json, err error) {

	os.Setenv("HTTP_PROXY", "http://127.0.0.1:6667")
	os.Setenv("HTTPS_PROXY", "https://127.0.0.1:6667")

	e.lastTimes++
	params = append(params, "api_key="+e.option.AccessKey)
	sort.Strings(params)
	params = append(params, "secret_key="+e.option.SecretKey)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(url, params)
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// GetAccount get the account detail of this exchange
func (e *OKEX) GetAccount() interface{} {
	json, err := e.getAuthJSON(e.host+"userinfo.do", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("GetAccount() error, the error number is %v", json.Get("error_code").MustInt())
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	return map[string]float64{
		"USDT":       conver.Float64Must(json.GetPath("info", "funds", "free", "usdt").Interface()),
		"FrozenUSDT": conver.Float64Must(json.GetPath("info", "funds", "freezed", "usdt").Interface()),
		"BTC":        conver.Float64Must(json.GetPath("info", "funds", "free", "btc").Interface()),
		"FrozenBTC":  conver.Float64Must(json.GetPath("info", "funds", "freezed", "btc").Interface()),
		"ETH":        conver.Float64Must(json.GetPath("info", "funds", "free", "eth").Interface()),
		"FrozenETH":  conver.Float64Must(json.GetPath("info", "funds", "freezed", "eth").Interface()),
		"EOS":        conver.Float64Must(json.GetPath("info", "funds", "free", "eos").Interface()),
		"FrozenEOS":  conver.Float64Must(json.GetPath("info", "funds", "freezed", "eos").Interface()),
		"ONT":        conver.Float64Must(json.GetPath("info", "funds", "free", "ont").Interface()),
		"FrozenONT":  conver.Float64Must(json.GetPath("info", "funds", "freezed", "ont").Interface()),
		"QTUM":       conver.Float64Must(json.GetPath("info", "funds", "free", "qtum").Interface()),
		"FrozenQTUM": conver.Float64Must(json.GetPath("info", "funds", "freezed", "qtum").Interface()),
	}
}

// Trade place an order
func (e *OKEX) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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

func (e *OKEX) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	params := []string{
		"symbol=" + e.stockTypeMap[stockType],
	}
	typeParam := "type=buy_market"
	amountParam := fmt.Sprintf("price=%f", amount)
	if price > 0 {
		typeParam = "type=buy"
		amountParam = fmt.Sprintf("amount=%f", amount)
		params = append(params, fmt.Sprintf("price=%f", price))
	}
	params = append(params, typeParam, amountParam)
	json, err := e.getAuthJSON(e.host+"trade.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("order_id").Interface())
}

func (e *OKEX) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	params := []string{
		"symbol=" + e.stockTypeMap[stockType],
		fmt.Sprintf("amount=%f", amount),
	}
	typeParam := "type=sell_market"
	if price > 0 {
		typeParam = "type=sell"
		params = append(params, fmt.Sprintf("price=%f", price))
	}
	params = append(params, typeParam)
	json, err := e.getAuthJSON(e.host+"trade.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("order_id").Interface())
}

// GetOrder get details of an order
func (e *OKEX) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType],
		"order_id=" + id,
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	ordersJSON := json.Get("orders")
	if len(ordersJSON.MustArray()) > 0 {
		orderJSON := ordersJSON.GetIndex(0)
		return Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		}
	}
	return false
}

// GetOrders get all unfilled orders
func (e *OKEX) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType],
		"order_id=-1",
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	orders := []Order{}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *OKEX) GetTrades(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType],
		"status=1",
		"current_page=1",
		"page_length=200",
	}
	json, err := e.getAuthJSON(e.host+"order_history.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	orders := []Order{}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder cancel an order
func (e *OKEX) CancelOrder(order Order) bool {
	params := []string{
		"symbol=" + e.stockTypeMap[order.StockType],
		"order_id=" + order.ID,
	}
	json, err := e.getAuthJSON(e.host+"cancel_order.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *OKEX) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("%vdepth.do?symbol=%v&size=%v", e.host, e.stockTypeMap[stockType], size))
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
func (e *OKEX) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *OKEX) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return false
	}
	if _, ok := e.recordsPeriodMap[period]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return false
	}
	size := 200
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("%vkline.do?symbol=%v&type=%v&size=%v", e.host, e.stockTypeMap[stockType], e.recordsPeriodMap[period], size))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	timeLast := int64(0)
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		recordTime := recordJSON.GetIndex(0).MustInt64() / 1000
		if recordTime > timeLast {
			recordsNew = append([]Record{{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[period][len(e.records[period])-1] = Record{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}
		} else {
			break
		}
	}
	e.records[period] = append(e.records[period], recordsNew...)
	if len(e.records[period]) > size {
		e.records[period] = e.records[period][len(e.records[period])-size : len(e.records[period])]
	}
	return e.records[period]
}
