package ZbAPI

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

// SHA1 加密
func digest() string {
	hash := sha1.New()
	hash.Write([]byte(Config.SECRET_KEY))
	return hex.EncodeToString(hash.Sum(nil))
}

// hmac MD5
func hmacSign(message string) string {
	hmac := hmac.New(md5.New, []byte(digest()))
	hmac.Write([]byte(message))
	return hex.EncodeToString(hmac.Sum(nil))
}

// 参数按照字母排序
func sortParams(params map[string]string) string {
	var buffer bytes.Buffer
	sortKey := make([]string, 0, len(params))

	for k := range params {
		sortKey = append(sortKey, k)
	}
	sort.Strings(sortKey)

	for _, k := range sortKey {
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(params[k])
		buffer.WriteString("&")
	}
	return strings.TrimSuffix(buffer.String(), "&")
}

type httpClient struct {
	*resty.Client
}

func (client *httpClient) handleQueryParams() {
	client.OnAfterResponse(func(client *resty.Client, req *resty.Response) error {
		for k := range client.QueryParam {
			delete(client.QueryParam, k)
		}
		return nil
	})

	client.OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
		client.SetQueryParams(map[string]string{
			"accesskey": Config.ACCESS_KEY,
			"reqTime":   strconv.FormatInt(time.Now().UnixNano()/1000000, 10),
		})
		return nil
	})

	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})
}
