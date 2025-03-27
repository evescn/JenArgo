package po

type ArgoCDUserInfo struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type ArgoCDRolloutInfo struct {
	ID int `json:"id"`
}
