package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/api/BinanceAPI"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// Binance the exchange struct of binance.com
type Binance struct {
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

// NewBinance create an exchange struct of Binance.com
func NewBinance(opt Option) Exchange {
	BinanceAPI.ACCESS_KEY = opt.AccessKey
	BinanceAPI.SECRET_KEY = opt.SecretKey
	return &Binance{
		stockTypeMap: map[string]string{
			"BTC/USDT":  "BTC",
			"ETH/USDT":  "ETH",
			"EOS/USDT":  "EOS",
			"ONT/USDT":  "ONT",
			"QTUM/USDT": "QTUM",
		},
		tradeTypeMap: map[string]string{
			"BUY":  constant.TradeTypeBuy,
			"SELL": constant.TradeTypeSell,
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
func (e *Binance) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *Binance) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *Binance) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *Binance) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *Binance) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *Binance) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *Binance) GetAccount() interface{} {
	accountsMap, err := BinanceAPI.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if _, ok := accountsMap["code"]; ok { //存在错误码
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", accountsMap["msg"].(string))
		return false
	}
	result := make(map[string]float64)
	balances := accountsMap["balances"].([]interface{})
	for _, n := range balances {
		//log.Println(n)
		b := n.(map[string]interface{}) //类型转换而已
		key := strings.ToUpper(b["asset"].(string))
		result[key] = conver.Float64Must(b["free"])
		result["Frozen"+key] = conver.Float64Must(b["locked"])
	}
	return result
}

// Trade place an order
func (e *Binance) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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

func (e *Binance) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	result, err := BinanceAPI.LimitBuy(conver.StringMust(amount), conver.StringMust(price), e.stockTypeMap[stockType]+"USDT")
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	orderId := conver.Int64Must(result["orderId"])
	if orderId <= 0 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", result["msg"].(string))
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return fmt.Sprint(orderId)
}

func (e *Binance) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	result, err := BinanceAPI.LimitSell(conver.StringMust(amount), conver.StringMust(price), e.stockTypeMap[stockType]+"USDT")
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	orderId := conver.Int64Must(result["orderId"])
	if orderId <= 0 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", result["msg"].(string))
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return fmt.Sprint(orderId)
}

// GetOrder get details of an order
func (e *Binance) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := BinanceAPI.GetOneOrder(id, e.stockTypeMap[stockType]+"USDT")
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if _, ok := result["code"]; ok { //存在错误码
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", result["msg"].(string))
		return false
	}
	return Order{
		ID:         id,
		Price:      conver.Float64Must(result["price"].(string)),
		Amount:     conver.Float64Must(result["origQty"].(string)),
		DealAmount: conver.Float64Must(result["executedQty"]),
		TradeType:  e.tradeTypeMap[result["side"].(string)],
		StockType:  stockType,
	}
}

// GetOrders get all unfilled orders
func (e *Binance) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	result, err := BinanceAPI.GetUnfinishOrders(e.stockTypeMap[stockType] + "USDT")
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	orders := []Order{}
	for _, n := range result {
		ord := n.(map[string]interface{})
		orders = append(orders, Order{
			ID:         fmt.Sprint(conver.Int64Must(ord["orderId"])),
			Price:      conver.Float64Must(ord["price"]),
			Amount:     conver.Float64Must(ord["origQty"]),
			DealAmount: conver.Float64Must(ord["executedQty"]),
			TradeType:  e.tradeTypeMap[ord["side"].(string)],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *Binance) GetTrades(stockType string) interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *Binance) CancelOrder(order Order) bool {
	ok, err := BinanceAPI.CancelOrder(order.ID, e.stockTypeMap[order.StockType]+"USDT")
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if ok {
		e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	}
	return ok
}

// getTicker get market ticker & depth
func (e *Binance) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	result, err := BinanceAPI.GetDepth(10, e.stockTypeMap[stockType]+"USDT")
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	if _, ok := result["code"]; ok { //存在错误码
		err = fmt.Errorf("GetTicker() error, %+v", result["msg"].(string))
		return
	}
	bids := result["bids"].([]interface{})
	asks := result["asks"].([]interface{})
	for _, bid := range bids {
		_bid := bid.([]interface{})
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  conver.Float64Must(_bid[0]),
			Amount: conver.Float64Must(_bid[1]),
		})
	}
	for _, ask := range asks {
		_ask := ask.([]interface{})
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  conver.Float64Must(_ask[0]),
			Amount: conver.Float64Must(_ask[1]),
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
func (e *Binance) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *Binance) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return nil
}
