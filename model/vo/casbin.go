package vo

type CasbinAuthRequest struct {
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

type CasbinAuthResponse struct {
	Code int         `json:"code"`           // 0 表示成功，其他值表示具体错误码
	Msg  string      `json:"msg"`            // 可读的提示信息
	Data interface{} `json:"data,omitempty"` // 成功时返回的数据
}
