package vo

type GitLabProjectListRequest struct {
	AppName string `json:"app_name" form:"app_name"`
	GroupId uint   `json:"group_id" form:"group_id"`
	Size    int    `form:"size" binding:"required"`
	Page    int    `form:"page" binding:"required"`
}

type GitLabProjectBranchRequest struct {
	AppId       uint `json:"app_id" form:"app_id"`
	BranchOrTag bool `json:"branch_or_tag" form:"branch_or_tag"`
}
