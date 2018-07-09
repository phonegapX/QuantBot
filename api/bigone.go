package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/api/BigoneAPI"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// BigOne the exchange struct of big.one
type BigOne struct {
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

var bo *BigoneAPI.Bigone

// NewBigOne create an exchange struct of big.one
func NewBigOne(opt Option) Exchange {
	bo = BigoneAPI.New(http.DefaultClient, opt.AccessKey, opt.SecretKey)
	//...
	return &BigOne{
		stockTypeMap: map[string]string{
			"BTC/USDT": "BTC-USDT",
			"ONE/USDT": "ONE-USDT",
			"EOS/USDT": "EOS-USDT",
			"ETH/USDT": "ETH-USDT",
			"BCH/USDT": "BCH-USDT",
			"EOS/ETH":  "EOS-ETH",
		},
		tradeTypeMap: map[string]string{
			"BID": constant.TradeTypeBuy,
			"ASK": constant.TradeTypeSell,
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
			"BTC/USDT": 0.001,
			"ONE/USDT": 0.001,
			"EOS/USDT": 0.001,
			"ETH/USDT": 0.001,
			"BCH/USDT": 0.001,
			"EOS/ETH":  0.001,
		},
		records: make(map[string][]Record),
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log print something to console
func (e *BigOne) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *BigOne) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *BigOne) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *BigOne) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *BigOne) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *BigOne) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *BigOne) GetAccount() interface{} {
	result, err := bo.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if len(result.Errors) > 0 {
		//log.Printf("response error : %v", result.Errors)
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", result.Errors[0].Message)
		return false
	}
	accInfo := make(map[string]float64)
	for _, v := range result.Data {
		available := conver.Float64Must(v.Balance)
		freez := conver.Float64Must(v.LockedBalance)
		if available != 0 {
			accInfo[strings.ToUpper(v.AssetID)] = available
		}
		if freez != 0 {
			accInfo["Frozen"+strings.ToUpper(v.AssetID)] = freez
		}
	}
	return accInfo
}

// Trade place an order
func (e *BigOne) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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

func (e *BigOne) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	result, err := bo.LimitBuy(conver.StringMust(amount), conver.StringMust(price), e.stockTypeMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if len(result.Errors) > 0 {
		//log.Printf("response error : %v", result.Errors)
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", result.Errors[0].Message)
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return result.Data.ID
}

func (e *BigOne) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	result, err := bo.LimitSell(conver.StringMust(amount), conver.StringMust(price), e.stockTypeMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if len(result.Errors) > 0 {
		//log.Printf("response error : %v", result.Errors)
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", result.Errors[0].Message)
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return result.Data.ID
}

// GetOrder get details of an order
func (e *BigOne) GetOrder(stockType, id string) interface{} {
	return nil
}

// GetOrders get all unfilled orders
func (e *BigOne) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := bo.GetUnfinishOrders(e.stockTypeMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if len(result.Errors) > 0 {
		//log.Printf("response error : %v", result.Errors)
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", result.Errors[0].Message)
		return false
	}
	orders := []Order{}
	for _, v := range result.Data.Edges {
		n := v.Node
		orders = append(orders, Order{
			ID:         n.ID,
			Price:      conver.Float64Must(n.Price),
			Amount:     conver.Float64Must(n.Amount),
			DealAmount: conver.Float64Must(n.FilledAmount),
			TradeType:  e.tradeTypeMap[n.Side],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *BigOne) GetTrades(stockType string) interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *BigOne) CancelOrder(order Order) bool {
	result, err := bo.CancelOrder(order.ID, e.stockTypeMap[order.StockType])
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if len(result.Errors) > 0 {
		//log.Printf("response error : %v", result.Errors)
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", result.Errors[0].Message)
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *BigOne) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	result, err := bo.GetDepth(e.stockTypeMap[stockType])
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	for _, bid := range result.Data.Bids {
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  conver.Float64Must(bid.Price),
			Amount: conver.Float64Must(bid.Amount),
		})
	}
	for _, ask := range result.Data.Asks {
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  conver.Float64Must(ask.Price),
			Amount: conver.Float64Must(ask.Amount),
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
func (e *BigOne) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *BigOne) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return nil
}
