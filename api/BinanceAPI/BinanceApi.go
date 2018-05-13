package BinanceAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	//"os"
	"strconv"
	"time"
)

const (
	API_BASE_URL = "https://api.binance.com/"
	API_V1       = API_BASE_URL + "api/v1/"
	API_V3       = API_BASE_URL + "api/v3/"

	TICKER_URI             = "ticker/24hr?symbol=%s"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "depth?symbol=%s&limit=%d"
	ACCOUNT_URI            = "account?"
	ORDER_URI              = "order?"
	UNFINISHED_ORDERS_INFO = "openOrders?"
)

var (
	ACCESS_KEY string       = ""
	SECRET_KEY string       = ""
	httpClient *http.Client = &http.Client{}
)

func init() {
	//os.Setenv("HTTP_PROXY", "http://127.0.0.1:6667")
	//os.Setenv("HTTPS_PROXY", "https://127.0.0.1:6667")
}

func buildParamsSigned(postForm *url.Values) error {
	postForm.Set("recvWindow", "6000000")
	tonce := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
	postForm.Set("timestamp", tonce)
	payload := postForm.Encode()
	sign, _ := GetParamHmacSHA256Sign(SECRET_KEY, payload)
	postForm.Set("signature", sign)
	return nil
}

func GetDepth(size int, symbol string) (map[string]interface{}, error) {
	if size > 100 {
		size = 100
	} else if size < 5 {
		size = 5
	}

	apiUrl := fmt.Sprintf(API_V1+DEPTH_URI, symbol, size)
	resp, err := HttpGet(httpClient, apiUrl)
	return resp, err
}

func GetAccount() (map[string]interface{}, error) {
	params := url.Values{}
	buildParamsSigned(&params)
	path := API_V3 + ACCOUNT_URI + params.Encode()
	respmap, err := HttpGet2(httpClient, path, map[string]string{"X-MBX-APIKEY": ACCESS_KEY})
	return respmap, err
}

func placeOrder(amount, price string, symbol string, orderType, orderSide string) (map[string]interface{}, error) {
	path := API_V3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", orderSide)
	params.Set("type", orderType)

	params.Set("quantity", amount)
	params.Set("type", "LIMIT")
	params.Set("timeInForce", "GTC")

	switch orderType {
	case "LIMIT":
		params.Set("price", price)
	}

	buildParamsSigned(&params)

	resp, err := HttpPostForm2(httpClient, path, params, map[string]string{"X-MBX-APIKEY": ACCESS_KEY})
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}

	return respmap, nil
}

func LimitBuy(amount, price string, symbol string) (map[string]interface{}, error) {
	return placeOrder(amount, price, symbol, "LIMIT", "BUY")
}

func LimitSell(amount, price string, symbol string) (map[string]interface{}, error) {
	return placeOrder(amount, price, symbol, "LIMIT", "SELL")
}

func MarketBuy(amount, price string, symbol string) (map[string]interface{}, error) {
	return placeOrder(amount, price, symbol, "MARKET", "BUY")
}

func MarketSell(amount, price string, symbol string) (map[string]interface{}, error) {
	return placeOrder(amount, price, symbol, "MARKET", "SELL")
}

func CancelOrder(orderId string, symbol string) (bool, error) {
	path := API_V3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("orderId", orderId)

	buildParamsSigned(&params)

	resp, err := HttpDeleteForm(httpClient, path, params, map[string]string{"X-MBX-APIKEY": ACCESS_KEY})

	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}

	orderIdCanceled := ToInt(respmap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func GetOneOrder(orderId string, symbol string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	if orderId != "" {
		params.Set("orderId", orderId)
	}
	params.Set("orderId", orderId)

	buildParamsSigned(&params)
	path := API_V3 + ORDER_URI + params.Encode()

	respmap, err := HttpGet2(httpClient, path, map[string]string{"X-MBX-APIKEY": ACCESS_KEY})
	return respmap, err
}

func GetUnfinishOrders(symbol string) ([]interface{}, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	buildParamsSigned(&params)
	path := API_V3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(httpClient, path, map[string]string{"X-MBX-APIKEY": ACCESS_KEY})
	return respmap, err
}
