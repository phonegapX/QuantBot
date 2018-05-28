package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/api/ZbAPI"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// Zb the exchange struct of zb.com
type Zb struct {
	stockTypeMap     map[string]string
	tradeTypeMap     map[int]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record
	logger           model.Logger
	option           Option

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewZb create an exchange struct of zb.com
func NewZb(opt Option) Exchange {
	ZbAPI.Config.ACCESS_KEY = opt.AccessKey
	ZbAPI.Config.SECRET_KEY = opt.SecretKey
	//...
	return &Zb{
		stockTypeMap: map[string]string{
			"BTC/USDT":  "btc_usdt",
			"ETH/USDT":  "eth_usdt",
			"EOS/USDT":  "eos_usdt",
			"LTC/USDT":  "ltc_usdt",
			"QTUM/USDT": "qtum_usdt",
		},
		tradeTypeMap: map[int]string{
			1: constant.TradeTypeBuy,
			0: constant.TradeTypeSell,
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
			"LTC/USDT":  0.001,
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
func (e *Zb) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *Zb) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *Zb) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *Zb) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *Zb) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *Zb) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *Zb) GetAccount() interface{} {
	accountInfo, err := ZbAPI.GetAccountInfo()
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if accountInfo.Code != 0 && accountInfo.Code != 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", accountInfo.Message)
		return false
	}
	result := make(map[string]float64)
	count := len(accountInfo.Result.Coins)
	for i := 0; i < count; i++ {
		coin := accountInfo.Result.Coins[i]
		freez := conver.Float64Must(coin.Freez)
		available := conver.Float64Must(coin.Available)
		if available != 0 {
			result[strings.ToUpper(coin.EnName)] = available
		}
		if freez != 0 {
			result["Frozen"+strings.ToUpper(coin.EnName)] = freez
		}
	}
	return result
}

// Trade place an order
func (e *Zb) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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

func (e *Zb) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	result, err := ZbAPI.CreateOrder(conver.StringMust(amount), e.stockTypeMap[stockType], "1", conver.StringMust(price))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if result.Code != 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", result.Message)
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return result.Id
}

func (e *Zb) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	result, err := ZbAPI.CreateOrder(conver.StringMust(amount), e.stockTypeMap[stockType], "0", conver.StringMust(price))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if result.Code != 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", result.Message)
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return result.Id
}

// GetOrder get details of an order
func (e *Zb) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := ZbAPI.GetOrder(id, e.stockTypeMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if result.Code != 0 && result.Code != 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", result.Message)
		return false
	}
	return Order{
		ID:         result.ID,
		Price:      result.Price,
		Amount:     result.TotalAmount,
		DealAmount: result.TradeAmount,
		TradeType:  e.tradeTypeMap[result.OrderType],
		StockType:  stockType,
	}
}

// GetOrders get all unfilled orders
func (e *Zb) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := ZbAPI.GetOrders(e.stockTypeMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	orders := []Order{}
	count := len(*result)
	for i := 0; i < count; i++ {
		orders = append(orders, Order{
			ID:         (*result)[i].ID,
			Price:      (*result)[i].Price,
			Amount:     (*result)[i].TotalAmount,
			DealAmount: (*result)[i].TradeAmount,
			TradeType:  e.tradeTypeMap[(*result)[i].OrderType],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *Zb) GetTrades(stockType string) interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *Zb) CancelOrder(order Order) bool {
	result, err := ZbAPI.CancelOrder(order.ID, e.stockTypeMap[order.StockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result.Code != 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", result.Message)
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *Zb) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	result, err := ZbAPI.GetDepth(e.stockTypeMap[stockType], "10")
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	count := len(result.Bids)
	for i := 0; i < count; i++ {
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  result.Bids[i][0],
			Amount: result.Bids[i][1],
		})
	}
	count = len(result.Asks)
	for i := 0; i < count; i++ {
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  result.Asks[i][0],
			Amount: result.Asks[i][1],
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
func (e *Zb) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *Zb) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return nil
}
