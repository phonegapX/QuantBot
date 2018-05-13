package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/phonegapX/QuantBot/api/HuobiProAPI/config"
	"github.com/phonegapX/QuantBot/api/HuobiProAPI/models"
	"github.com/phonegapX/QuantBot/api/HuobiProAPI/untils"
)

// 批量操作的API下个版本再封装

//------------------------------------------------------------------------------------------
// 交易API

// 获取K线数据
// strSymbol: 交易对, btcusdt, bccbtc......
// strPeriod: K线类型, 1min, 5min, 15min......
// nSize: 获取数量, [1-2000]
// return: KLineReturn 对象
func GetKLine(strSymbol, strPeriod string, nSize int) (r models.KLineReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["period"] = strPeriod
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/kline"
	strUrl := config.MARKET_URL + strRequestUrl

	jsonKLineReturn := untils.HttpGetRequest(strUrl, mapParams)
	err = json.Unmarshal([]byte(jsonKLineReturn), &r)

	return
}

// 获取聚合行情
// strSymbol: 交易对, btcusdt, bccbtc......
// return: TickReturn对象
func GetTicker(strSymbol string) (r models.TickerReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail/merged"
	strUrl := config.MARKET_URL + strRequestUrl

	jsonTickReturn := untils.HttpGetRequest(strUrl, mapParams)
	err = json.Unmarshal([]byte(jsonTickReturn), &r)

	return
}

// 获取交易深度信息
// strSymbol: 交易对, btcusdt, bccbtc......
// strType: Depth类型, step0、step1......stpe5 (合并深度0-5, 0时不合并)
// return: MarketDepthReturn对象
func GetMarketDepth(strSymbol, strType string) (r models.MarketDepthReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["type"] = strType

	strRequestUrl := "/market/depth"
	strUrl := config.MARKET_URL + strRequestUrl

	jsonMarketDepthReturn := untils.HttpGetRequest(strUrl, mapParams)
	err = json.Unmarshal([]byte(jsonMarketDepthReturn), &r)

	return
}

// 获取交易细节信息
// strSymbol: 交易对, btcusdt, bccbtc......
// return: TradeDetailReturn对象
func GetTradeDetail(strSymbol string) (r models.TradeDetailReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/trade"
	strUrl := config.MARKET_URL + strRequestUrl

	jsonTradeDetailReturn := untils.HttpGetRequest(strUrl, mapParams)
	err = json.Unmarshal([]byte(jsonTradeDetailReturn), &r)

	return
}

// 批量获取最近的交易记录
// strSymbol: 交易对, btcusdt, bccbtc......
// nSize: 获取交易记录的数量, 范围1-2000
// return: TradeReturn对象
func GetTrade(strSymbol string, nSize int) (r models.TradeReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/trade"
	strUrl := config.MARKET_URL + strRequestUrl

	jsonTradeReturn := untils.HttpGetRequest(strUrl, mapParams)
	err = json.Unmarshal([]byte(jsonTradeReturn), &r)

	return
}

// 获取Market Detail 24小时成交量数据
// strSymbol: 交易对, btcusdt, bccbtc......
// return: MarketDetailReturn对象
func GetMarketDetail(strSymbol string) (r models.MarketDetailReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail"
	strUrl := config.MARKET_URL + strRequestUrl

	jsonMarketDetailReturn := untils.HttpGetRequest(strUrl, mapParams)
	err = json.Unmarshal([]byte(jsonMarketDetailReturn), &r)

	return
}

//------------------------------------------------------------------------------------------
// 公共API

// 查询系统支持的所有交易及精度
// return: SymbolsReturn对象
func GetSymbols() (r models.SymbolsReturn, err error) {
	strRequestUrl := "/v1/common/symbols"
	strUrl := config.TRADE_URL + strRequestUrl

	jsonSymbolsReturn := untils.HttpGetRequest(strUrl, nil)
	err = json.Unmarshal([]byte(jsonSymbolsReturn), &r)

	return
}

// 查询系统支持的所有币种
// return: CurrencysReturn对象
func GetCurrencys() (r models.CurrencysReturn, err error) {
	strRequestUrl := "/v1/common/currencys"
	strUrl := config.TRADE_URL + strRequestUrl

	jsonCurrencysReturn := untils.HttpGetRequest(strUrl, nil)
	err = json.Unmarshal([]byte(jsonCurrencysReturn), &r)

	return
}

// 查询系统当前时间戳
// return: TimestampReturn对象
func GetTimestamp() (r models.TimestampReturn, err error) {
	strRequest := "/v1/common/timestamp"
	strUrl := config.TRADE_URL + strRequest

	jsonTimestampReturn := untils.HttpGetRequest(strUrl, nil)
	err = json.Unmarshal([]byte(jsonTimestampReturn), &r)

	return
}

//------------------------------------------------------------------------------------------
// 用户资产API

// 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func GetAccounts() (r models.AccountsReturn, err error) {
	strRequest := "/v1/account/accounts"

	jsonAccountsReturn := untils.ApiKeyGet(make(map[string]string), strRequest)
	err = json.Unmarshal([]byte(jsonAccountsReturn), &r)

	return
}

// 根据账户ID查询账户余额
// nAccountID: 账户ID, 不知道的话可以通过GetAccounts()获取, 可以只现货账户, C2C账户, 期货账户
// return: BalanceReturn对象
func GetAccountBalance(strAccountID string) (r models.BalanceReturn, err error) {
	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strAccountID)

	jsonBanlanceReturn := untils.ApiKeyGet(make(map[string]string), strRequest)
	err = json.Unmarshal([]byte(jsonBanlanceReturn), &r)

	return
}

//------------------------------------------------------------------------------------------
// 交易API

// 下单
// params: 下单信息
// return: PlaceReturn对象
func Place(params models.PlaceRequestParams) (r models.PlaceReturn, err error) {
	mapParams := make(map[string]string)
	mapParams["account-id"] = params.AccountID
	mapParams["amount"] = params.Amount
	if 0 < len(params.Price) {
		mapParams["price"] = params.Price
	}
	if 0 < len(params.Source) {
		mapParams["source"] = params.Source
	}
	mapParams["symbol"] = params.Symbol
	mapParams["type"] = params.Type

	strRequest := "/v1/order/orders/place"

	jsonPlaceReturn := untils.ApiKeyPost(mapParams, strRequest)
	err = json.Unmarshal([]byte(jsonPlaceReturn), &r)

	return
}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func SubmitCancel(strOrderID string) (r models.PlaceReturn, err error) {
	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", strOrderID)

	jsonPlaceReturn := untils.ApiKeyPost(make(map[string]string), strRequest)
	err = json.Unmarshal([]byte(jsonPlaceReturn), &r)

	return
}

// 根据订单ID查询订单详情
func GetOrderDetail(strOrderID string) (r models.OrderDetailReturn, err error) {
	strRequest := fmt.Sprintf("/v1/order/orders/%s", strOrderID)

	jsonOrderReturn := untils.ApiKeyGet(make(map[string]string), strRequest)
	err = json.Unmarshal([]byte(jsonOrderReturn), &r)

	return
}

// 列出当前所有挂单
func GetOrders(strSymbol string) (r models.OrdersReturn, err error) {
	//pre-submitted 准备提交, submitted 已提交, partial-filled 部分成交, partial-canceled 部分成交撤销, filled 完全成交, canceled 已撤销
	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["states"] = "submitted,partial-filled"

	strRequest := "/v1/order/orders"

	jsonOrdersReturn := untils.ApiKeyGet(mapParams, strRequest)
	err = json.Unmarshal([]byte(jsonOrdersReturn), &r)

	return
}
