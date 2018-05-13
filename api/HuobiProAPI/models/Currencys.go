package models

type CurrencysReturn struct {
	Status  string   `json:"status"` // 请求状态
	Data    []string `json:"data"`   // 系统支持的所有币种
	ErrCode string   `json:"err-code"`
	ErrMsg  string   `json:"err-msg"`
}
