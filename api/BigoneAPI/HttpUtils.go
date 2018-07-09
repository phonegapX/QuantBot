package BigoneAPI

//http request 工具函数
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func NewHttpRequest(client *http.Client, reqType string, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	if requstHeaders != nil {
		for k, v := range requstHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	return bodyData, nil
}

func HttpGet(client *http.Client, reqUrl string, headers map[string]string, result interface{}) error {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewHttpRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respData, result)
	if err != nil {
		log.Printf("HttpGet4 - json.Unmarshal failed : %v, resp %s", err, string(respData))
		return err
	}

	return nil
}

func HttpPostForm(client *http.Client, reqUrl string, postData map[string]string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/json"
	data, _ := json.Marshal(postData)
	return NewHttpRequest(client, "POST", reqUrl, string(data), headers)
}
