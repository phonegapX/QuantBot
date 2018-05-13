package untils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	//"os"
	"sort"
	"strings"
	"time"

	"github.com/phonegapX/QuantBot/api/HuobiProAPI/config"
	//"golang.org/x/net/proxy"
)

// Http Get请求基础函数, 通过封装Go语言Http请求, 支持火币网REST API的HTTP Get请求
// strUrl: 请求的URL
// strParams: string类型的请求参数, user=lxz&pwd=lxz
// return: 请求结果
func HttpGetRequest(strUrl string, mapParams map[string]string) string {

	//=============================================================
	// create a socks5 dialer
	//dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:6666", nil, proxy.Direct)
	//if err != nil {
	//	return err.Error()
	//}
	// setup a http client
	//httpTransport := &http.Transport{}
	//httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	//httpTransport.Dial = dialer.Dial

	//==========================================================

	//os.Setenv("HTTP_PROXY", "http://127.0.0.1:6667")
	//os.Setenv("HTTPS_PROXY", "https://127.0.0.1:6667")

	//==========================================================
	//
	httpClient := &http.Client{}

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := Map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	// 发出请求
	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if nil != err {
		return err.Error()
	}

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

// Http POST请求基础函数, 通过封装Go语言Http请求, 支持火币网REST API的HTTP POST请求
// strUrl: 请求的URL
// mapParams: map类型的请求参数
// return: 请求结果
func HttpPostRequest(strUrl string, mapParams map[string]string) string {

	//=============================================================
	// create a socks5 dialer
	//dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:6666", nil, proxy.Direct)
	//if err != nil {
	//	return err.Error()
	//}
	// setup a http client
	//httpTransport := &http.Transport{}
	//httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	//httpTransport.Dial = dialer.Dial

	//==========================================================

	//os.Setenv("HTTP_PROXY", "http://127.0.0.1:6667")
	//os.Setenv("HTTPS_PROXY", "https://127.0.0.1:6667")

	//==========================================================
	//
	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")

	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if nil != err {
		return err.Error()
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

// 进行签名后的HTTP GET请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	strMethod := "GET"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams["AccessKeyId"] = config.ACCESS_KEY
	mapParams["SignatureMethod"] = "HmacSHA256"
	mapParams["SignatureVersion"] = "2"
	mapParams["Timestamp"] = timestamp

	hostName := "api.huobi.pro"
	mapParams["Signature"] = CreateSign(mapParams, strMethod, hostName, strRequestPath, config.SECRET_KEY)

	strUrl := config.TRADE_URL + strRequestPath
	return HttpGetRequest(strUrl, MapValueEncodeURI(mapParams))
}

// 进行签名后的HTTP POST请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
	strMethod := "POST"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = config.ACCESS_KEY
	mapParams2Sign["SignatureMethod"] = "HmacSHA256"
	mapParams2Sign["SignatureVersion"] = "2"
	mapParams2Sign["Timestamp"] = timestamp

	hostName := "api.huobi.pro"

	mapParams2Sign["Signature"] = CreateSign(mapParams2Sign, strMethod, hostName, strRequestPath, config.SECRET_KEY)
	strUrl := config.TRADE_URL + strRequestPath + "?" + Map2UrlQuery(MapValueEncodeURI(mapParams2Sign))

	return HttpPostRequest(strUrl, mapParams)
}

// 构造签名
// mapParams: 送进来参与签名的参数, Map类型
// strMethod: 请求的方法 GET, POST......
// strHostUrl: 请求的主机
// strRequestPath: 请求的路由路径
// strSecretKey: 进行签名的密钥
func CreateSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	// 参数处理, 按API要求, 参数名应按ASCII码进行排序(使用UTF-8编码, 其进行URI编码, 16进制字符必须大写)
	sortedParams := MapSortByKey(mapParams)
	encodeParams := MapValueEncodeURI(sortedParams)
	strParams := Map2UrlQuery(encodeParams)

	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + strParams

	return ComputeHmac256(strPayload, strSecretKey)
}

// 对Map按着ASCII码进行排序
// mapValue: 需要进行排序的map
// return: 排序后的map
func MapSortByKey(mapValue map[string]string) map[string]string {
	var keys []string
	for key := range mapValue {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	mapReturn := make(map[string]string)
	for _, key := range keys {
		mapReturn[key] = mapValue[key]
	}

	return mapReturn
}

// 对Map的值进行URI编码
// mapParams: 需要进行URI编码的map
// return: 编码后的map
func MapValueEncodeURI(mapValue map[string]string) map[string]string {
	for key, value := range mapValue {
		valueEncodeURI := url.QueryEscape(value)
		mapValue[key] = valueEncodeURI
	}

	return mapValue
}

// 将map格式的请求参数转换为字符串格式的
// mapParams: map格式的参数键值对
// return: 查询字符串
func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	for key, value := range mapParams {
		strParams += (key + "=" + value + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

// HMAC SHA256加密
// strMessage: 需要加密的信息
// strSecret: 密钥
// return: BASE64编码的密文
func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
