package vo

import "JenArgo/model/bo"

type ArgoCDApp struct {
	Name      string `json:"name" form:"name"`
	Namespace string `json:"namespace" form:"namespace"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

type ArgoCDImage struct {
	Name string `json:"name" form:"name"`
}

type ArgoCDRollback struct {
	Name       string `json:"name" form:"name"`
	RollbackID int    `json:"rollback_id" form:"rollback_id"`
}
type ArgoCDImageRequest struct {
	AppName    string  `json:"appName"`
	AppProject string  `json:"appProject"`
	AppID      int     `json:"appID"`
	Source     *Source `json:"source"`
}

type Source struct {
	AppName string `json:"appName"`
	Helm    struct {
		ReleaseName string   `json:"releaseName"`
		ValueFiles  []string `json:"valueFiles"`
	} `json:"helm"`
	Path           string `json:"path"`
	RepoURL        string `json:"repoURL"`
	TargetRevision string `json:"targetRevision"`
}

type ArgoCDImageResponse struct {
	Images []*bo.ImageInfo `json:"images"`
}
