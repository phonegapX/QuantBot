package ZbAPI

import (
	"os"

	"github.com/go-resty/resty"
)

type configure struct {
	ACCESS_KEY string
	SECRET_KEY string
	dataURL    string
	tradeURL   string
}

// Config 中币接口配置信息
var (
	Config                  configure
	dataClient, tradeClient httpClient
)

func init() {

	os.Setenv("HTTP_PROXY", "http://127.0.0.1:6667")
	os.Setenv("HTTPS_PROXY", "https://127.0.0.1:6667")

	Config.ACCESS_KEY = ""
	Config.SECRET_KEY = ""
	Config.dataURL = "http://api.zb.com/data/v1/"
	Config.tradeURL = "https://trade.zb.com/api/"

	c1 := resty.New().SetDebug(false).SetHostURL(Config.dataURL)
	c2 := resty.New().SetDebug(false).SetHostURL(Config.tradeURL)
	dataClient = httpClient{c1}
	tradeClient = httpClient{c2}

	dataClient.handleQueryParams()
	tradeClient.handleQueryParams()
}
