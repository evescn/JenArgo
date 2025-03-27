package po

type GitLabProjectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	NameSpaceId uint   `json:"namespace_id"`
	Visibility  string `json:"visibility"`
	ImportUrl   string `json:"import_url"`
}
