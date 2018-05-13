package models

type AccountsData struct {
	ID     int64  `json:"id"`      // Account ID
	Type   string `json:"type"`    // 账户类型, spot: 现货账户
	State  string `json:"state"`   // 账户状态, working: 正常, lock: 账户被锁定
	UserID int64  `json:"user-id"` // 用户ID
}

type AccountsReturn struct {
	Status  string         `json:"status"` // 请求状态
	Data    []AccountsData `json:"data"`   // 用户数据
	ErrCode string         `json:"err-code"`
	ErrMsg  string         `json:"err-msg"`
}
