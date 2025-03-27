package vo

type JenArgoCiCd struct {
	En       string `json:"en"`
	AppName  string `json:"app_name"`
	RepoName string `json:"repo_name"`
	Builder  string `json:"builder"`
}

type JenArgoUpdateCiCd struct {
	En           string `json:"en"`
	AppName      string `json:"app_name"`
	RepoName     string `json:"repo_name"`
	Branch       string `json:"branch"`
	CodeCheck    int    `json:"code_check"`
	BuildStatus  int    `json:"build_status"`
	DeployStatus int    `json:"deploy_status"`
}
