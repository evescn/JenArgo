package bo

type GroupsInfo []struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProjectInfo []struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	NameSpace   *Namespace `json:"namespace"`
	WebUrl      string     `json:"web_url"`
	Description string     `json:"description"`
	CreatedAt   string     `json:"created_at"`
}

type BranchInfo []struct {
	Name string `json:"name"`
}

type Namespace struct {
	Name string `json:"name"`
}

type GitLabRequest struct {
	GroupName   string `json:"group_name" form:"group_name"`
	ProjectName string `json:"project_name" form:"project_name"`
	Visibility  string `json:"visibility" form:"visibility"`
	Description string `json:"desc" form:"desc"`
	GroupId     uint   `json:"group_id" form:"group_id"`
}
