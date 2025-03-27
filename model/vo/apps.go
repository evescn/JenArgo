package vo

type AppAddRequest struct {
	GroupName   string `json:"group_name" form:"group_name"`
	ProjectName string `json:"project_name" form:"project_name"`
	GroupId     uint   `json:"group_id" form:"group_id"`
	Visibility  string `json:"visibility" form:"visibility"`
	Description string `json:"desc" form:"desc"`
	HasJenkins  bool   `json:"has_jenkins" form:"has_jenkins"`
}
