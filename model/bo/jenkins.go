package bo

type Jenkins struct {
	GroupName   string `json:"group_name"`
	JobName     string `json:"job_name"`
	CopyJobName string `json:"copy_job_name"`
}
