package ZbAPI

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// 市场深度
// depth("depth", "btc_usdt", "20")
func depth(api, market, size string) (*respDepth, error) {
	resp, err := dataClient.SetQueryParams(map[string]string{
		"market": market,
		"size":   size,
	}).R().Get(api)
	if err != nil {
		return nil, err
	}
	var (
		res          respDepth
		mapInterface map[string]interface{}
	)
	err = json.Unmarshal(resp.Body(), &mapInterface)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(mapInterface, &res)
	return &res, err
}

func GetDepth(market, size string) (*respDepth, error) {
	return depth("depth", market, size)
}

// 行情
// getTicker("ticker", "btc_usdt")
func getTicker(api, market string) *respTicker {
	resp, _ := dataClient.SetQueryParams(map[string]string{
		"market": market,
	}).R().Get(api)
	var res respTicker
	json.Unmarshal(resp.Body(), &res)
	return &res
}

// K线
// kline("kline", "btc_usdt", "1min", "10")
func kline(api, market, timeType, size string) *respKline {
	resp, _ := dataClient.SetQueryParams(map[string]string{
		"market": market,
		"type":   timeType,
		"size":   size,
	}).R().Get(api)
	var (
		res          respKline
		mapInterface map[string]interface{}
	)
	json.Unmarshal(resp.Body(), &res)
	mapstructure.Decode(mapInterface, &res)
	return &res
}

// 历史成交
// trades("trades", "btc_usdt")
func trades(api, market string) *respTrades {
	resp, _ := dataClient.SetQueryParams(map[string]string{
		"market": market,
	}).R().Get(api)
	var res respTrades
	json.Unmarshal(resp.Body(), &res)
	return &res
}

// 获取用户信息
func accountInfo(api, sign string) (*respAccountInfo, error) {
	resp, err := tradeClient.SetQueryParams(map[string]string{
		"method": api,
		"sign":   sign,
	}).R().Get(api)
	if err != nil {
		return nil, err
	}
	var (
		res          respAccountInfo
		mapInterface map[string]interface{}
	)
	err = json.Unmarshal(resp.Body(), &mapInterface)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(mapInterface, &res)
	return &res, err
}

func GetAccountInfo() (*respAccountInfo, error) {
	params := map[string]string{
		"accesskey": Config.ACCESS_KEY,
		"method":    "getAccountInfo",
	}
	sorted := sortParams(params)
	sign := hmacSign(sorted)
	return accountInfo("getAccountInfo", sign)
}

// 委托下单
func createOrder(api, amount, currency, tradeType, price, sign string) (*respOrder, error) {
	resp, err := tradeClient.SetQueryParams(map[string]string{
		"amount":    amount,
		"currency":  currency,
		"method":    api,
		"price":     price,
		"tradeType": tradeType,
		"sign":      sign,
	}).R().Get(api)
	if err != nil {
		return nil, err
	}
	var res respOrder
	err = json.Unmarshal(resp.Body(), &res)
	return &res, err
}

func CreateOrder(amount, currency, tradeType, price string) (*respOrder, error) {
	createOrderParams := map[string]string{
		"accesskey": Config.ACCESS_KEY,
		"amount":    amount,
		"currency":  currency,
		"price":     price,
		"tradeType": tradeType,
		"method":    "order",
	}
	createOrderSorted := sortParams(createOrderParams)
	createOrderSign := hmacSign(createOrderSorted)
	return createOrder("order", amount, currency, tradeType, price, createOrderSign)
}

// 获取委托买单和卖单
func getOrders(api, currency, sign string) (*respOrders, error) {
	resp, err := tradeClient.SetQueryParams(map[string]string{
		"currency":  currency,
		"method":    api,
		"pageIndex": "1",
		"pageSize":  "10",
		"sign":      sign,
	}).R().Get(api)
	if err != nil {
		return nil, err
	}
	if resp.Body()[0] == '{' {
		var res respSimple
		err = json.Unmarshal(resp.Body(), &res)
		if err != nil {
			return nil, err
		}
		err = fmt.Errorf("%+v", res.Message)
		return nil, err
	}

	var res respOrders
	err = json.Unmarshal(resp.Body(), &res)
	return &res, err
}

func GetOrders(currency string) (*respOrders, error) {
	orderParams := map[string]string{
		"accesskey": Config.ACCESS_KEY,
		"currency":  currency,
		"method":    "getUnfinishedOrdersIgnoreTradeType",
		"pageIndex": "1",
		"pageSize":  "10",
	}
	orderSorted := sortParams(orderParams)
	orderSign := hmacSign(orderSorted)
	return getOrders("getUnfinishedOrdersIgnoreTradeType", currency, orderSign)
}

// 取消委托
func cancelOrder(api, id, currency, sign string) (*respSimple, error) {
	resp, err := tradeClient.SetQueryParams(map[string]string{
		"currency": currency,
		"method":   api,
		"id":       id,
		"sign":     sign,
	}).R().Get(api)
	if err != nil {
		return nil, err
	}
	var res respSimple
	err = json.Unmarshal(resp.Body(), &res)
	return &res, err
}

func CancelOrder(id, currency string) (*respSimple, error) {
	cancelParams := map[string]string{
		"accesskey": Config.ACCESS_KEY,
		"currency":  currency,
		"id":        id,
		"method":    "cancelOrder",
	}
	cancelSorted := sortParams(cancelParams)
	cancelSign := hmacSign(cancelSorted)
	return cancelOrder("cancelOrder", id, currency, cancelSign)
}

// 获取委托订单
func getOrder(api, id, currency, sign string) (*order, error) {
	resp, err := tradeClient.SetQueryParams(map[string]string{
		"currency": currency,
		"method":   api,
		"id":       id,
		"sign":     sign,
	}).R().Get(api)
	if err != nil {
		return nil, err
	}
	var res order
	err = json.Unmarshal(resp.Body(), &res)
	return &res, err
}

func GetOrder(id, currency string) (*order, error) {
	orderParams := map[string]string{
		"accesskey": Config.ACCESS_KEY,
		"currency":  currency,
		"id":        id,
		"method":    "getOrder",
	}
	orderSorted := sortParams(orderParams)
	orderSign := hmacSign(orderSorted)
	return getOrder("getOrder", id, currency, orderSign)
}
