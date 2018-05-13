package models

// 子账户结构
type SubAccount struct {
	Currency string `json:"currency"` // 币种
	Balance  string `json:"balance"`  // 结余
	Type     string `json:"type"`     // 类型, trade: 交易余额, frozen: 冻结余额
}

type Balance struct {
	ID     int64        `json:"id"`    // 账户ID
	State  string       `json:"state"` // 账户状态, working: 正常, lock: 账户被锁定
	Type   string       `json:"type"`  // 账户类型, spot: 现货账户
	List   []SubAccount `json:"list"`  // 子账户数组
	UserID int64        `json:"user-id"`
}

type BalanceReturn struct {
	Status  string  `json:"status"` // 请求状态
	Data    Balance `json:"data"`   // 账户余额
	ErrCode string  `json:"err-code"`
	ErrMsg  string  `json:"err-msg"`
}
