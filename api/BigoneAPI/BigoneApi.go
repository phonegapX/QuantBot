package BigoneAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nubo/jwt"
)

const (
	API_BASE_URL = "https://big.one/api/v2"
	TICKER_URI   = API_BASE_URL + "/markets/%s/ticker"
	DEPTH_URI    = API_BASE_URL + "/markets/%s/depth"
	ACCOUNT_URI  = API_BASE_URL + "/viewer/accounts"
	ORDERS_URI   = API_BASE_URL + "/viewer/orders"
)

type Bigone struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, api_key, secret_key string) *Bigone {
	return &Bigone{api_key, secret_key, client}
}

type TickerResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		Ask struct {
			Amount string `json:"amount"`
			Price  string `json:"price"`
		} `json:"ask"`
		Bid struct {
			Amount string `json:"amount"`
			Price  string `json:"price"`
		} `json:"bid"`
		Close           string `json:"close"`
		DailyChange     string `json:"daily_change"`
		DailyChangePerc string `json:"daily_change_perc"`
		High            string `json:"high"`
		Low             string `json:"low"`
		MarketID        string `json:"market_id"`
		MarketUUID      string `json:"market_uuid"`
		Open            string `json:"open"`
		Volume          string `json:"volume"`
	} `json:"data"`
}

func (bo *Bigone) GetTicker(currencyPair string) (*TickerResp, error) {
	var resp TickerResp
	tickerURI := fmt.Sprintf(TICKER_URI, currencyPair)
	err := HttpGet(bo.httpClient, tickerURI, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type PlaceOrderResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		Amount       string `json:"amount"`
		AvgDealPrice string `json:"avg_deal_price"`
		FilledAmount string `json:"filled_amount"`
		ID           string `json:"id"`
		InsertedAt   string `json:"inserted_at"`
		MarketID     string `json:"market_id"`
		MarketUUID   string `json:"market_uuid"`
		Price        string `json:"price"`
		Side         string `json:"side"`
		State        string `json:"state"`
		UpdatedAt    string `json:"updated_at"`
	} `json:"data"`
}

func (bo *Bigone) placeOrder(amount, price string, currencyPair string, orderType, orderSide string) (*PlaceOrderResp, error) {
	path := ORDERS_URI
	params := make(map[string]string)
	params["market_id"] = currencyPair
	params["side"] = orderSide
	params["amount"] = amount
	params["price"] = price

	var resp PlaceOrderResp
	buf, err := HttpPostForm(bo.httpClient, path, params, bo.privateHeader())
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf, &resp); nil != err {
		log.Printf("buf : %s", string(buf))
		log.Printf("placeOrder - json.Unmarshal failed : %v", err)
		return nil, err
	}

	return &resp, nil
}

func (bo *Bigone) LimitBuy(amount, price string, currencyPair string) (*PlaceOrderResp, error) {
	return bo.placeOrder(amount, price, currencyPair, "LIMIT", "BID")
}

func (bo *Bigone) LimitSell(amount, price string, currencyPair string) (*PlaceOrderResp, error) {
	return bo.placeOrder(amount, price, currencyPair, "LIMIT", "ASK")
}

func (bo *Bigone) MarketBuy(amount, price string, currencyPair string) (*PlaceOrderResp, error) {
	panic("not implements")
}

func (bo *Bigone) MarketSell(amount, price string, currencyPair string) (*PlaceOrderResp, error) {
	panic("not implements")
}

func (bo *Bigone) privateHeader() map[string]string {
	claims := jwt.ClaimSet{
		"type":  "OpenAPI",
		"sub":   bo.accessKey,
		"nonce": time.Now().UnixNano(),
	}
	token, err := claims.Sign(bo.secretKey)
	if nil != err {
		log.Printf("privateHeader - cliam.Sign failed : %v", err)
		return nil
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

type OrderListResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		Edges []struct {
			Cursor string `json:"cursor"`
			Node   struct {
				Amount       string `json:"amount"`
				AvgDealPrice string `json:"avg_deal_price"`
				FilledAmount string `json:"filled_amount"`
				ID           string `json:"id"`
				InsertedAt   string `json:"inserted_at"`
				MarketID     string `json:"market_id"`
				MarketUUID   string `json:"market_uuid"`
				Price        string `json:"price"`
				Side         string `json:"side"`
				State        string `json:"state"`
				UpdatedAt    string `json:"updated_at"`
			} `json:"node"`
		} `json:"edges"`
		PageInfo struct {
			EndCursor       string `json:"end_cursor"`
			HasNextPage     bool   `json:"has_next_page"`
			HasPreviousPage bool   `json:"has_previous_page"`
			StartCursor     string `json:"start_cursor"`
		} `json:"page_info"`
	} `json:"data"`
}

func (bo *Bigone) getOrdersList(currencyPair string, size int, tpy int) (*OrderListResp, error) {
	apiURL := ""
	apiURL = fmt.Sprintf("%s?market_id=%s", ORDERS_URI, currencyPair)

	if tpy == 0 {
		apiURL += "&state=FILLED"
	} else {
		apiURL += "&state=PENDING"
	}

	var resp OrderListResp
	err := HttpGet(bo.httpClient, apiURL, bo.privateHeader(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (bo *Bigone) GetUnfinishOrders(currencyPair string) (*OrderListResp, error) {
	return bo.getOrdersList(currencyPair, -1, 1)
}

func (bo *Bigone) GetOrderHistorys(currencyPair string) (*OrderListResp, error) {
	return bo.getOrdersList(currencyPair, -1, 0)
}

type CancelOrderResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		ID           string `json:"id"`
		MarketUUID   string `json:"market_uuid"`
		Price        string `json:"price"`
		Amount       string `json:"amount"`
		FilledAmount string `json:"filled_amount"`
		AvgDealPrice string `json:"avg_deal_price"`
		Side         string `json:"side"`
		State        string `json:"state"`
	}
}

func (bo *Bigone) CancelOrder(orderId string, currencyPair string) (*CancelOrderResp, error) {
	path := ORDERS_URI + "/" + orderId + "/cancel"
	params := make(map[string]string)
	params["order_id"] = orderId

	buf, err := HttpPostForm(bo.httpClient, path, params, bo.privateHeader())
	if err != nil {
		return nil, err
	}

	var resp CancelOrderResp
	if err = json.Unmarshal(buf, &resp); nil != err {
		log.Printf("CancelOrder - json.Unmarshal failed : %v", err)
		return nil, err
	}

	return &resp, nil
}

type AccountResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data []struct {
		AssetID       string `json:"asset_id"`
		AssetUUID     string `json:"asset_uuid"`
		Balance       string `json:"balance"`
		LockedBalance string `json:"locked_balance"`
	} `json:"data"`
}

func (bo *Bigone) GetAccount() (*AccountResp, error) {
	var resp AccountResp
	apiUrl := ACCOUNT_URI

	err := HttpGet(bo.httpClient, apiUrl, bo.privateHeader(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type DepthResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`

	Data struct {
		MarketID string `json:"market_id"`
		Bids     []struct {
			Price      string `json:"price"`
			OrderCount int    `json:"order_count"`
			Amount     string `json:"amount"`
		} `json:"bids"`
		Asks []struct {
			Price      string `json:"price"`
			OrderCount int    `json:"order_count"`
			Amount     string `json:"amount"`
		} `json:"asks"`
	}
}

func (bo *Bigone) GetDepth(currencyPair string) (*DepthResp, error) {
	var resp DepthResp
	apiURL := fmt.Sprintf(DEPTH_URI, currencyPair)
	err := HttpGet(bo.httpClient, apiURL, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
