package vo

import "JenArgo/model/po"

type DeploysListRequest struct {
	En       string `form:"en"`
	AppName  string `form:"app_name"`
	RepoName string `form:"repo_name"`
	Page     int    `form:"page"`
	Size     int    `form:"size"`
}

type DeployRequest struct {
	*po.Deploy
}

func (d *DeployRequest) ToDeploy() *po.Deploy {
	return d.Deploy // 返回一个新的 po.Deploy
}

type DeployDelRequest struct {
	ID int64 `json:"id"`
}

type DeploysListResponse struct {
	Items []*po.Deploy `json:"items"`
	Total int          `json:"total"`
}
