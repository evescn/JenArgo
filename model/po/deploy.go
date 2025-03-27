package po

import "time"

type Deploy struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	AppName          string `json:"app_name"`
	RepoName         string `json:"repo_name"`
	En               string `json:"en"`
	Branch           string `json:"branch"`
	CodeCheck        int    `json:"code_check" gorm:"column:code_check"`
	Tag              bool   `json:"tag"`
	Status           int    `json:"status"`
	StartTime        string `json:"start_time" gorm:"column:start_time"`
	Duration         string `json:"duration"`
	BuildStatus      int    `json:"build_status" gorm:"column:build_status"`
	DeployStatus     int    `json:"deploy_status" gorm:"column:deploy_status"`
	Builder          string `json:"builder"`
	BuildUrl         string `json:"build_url" gorm:"column:build_url"`
	HasScheduledTask bool   `json:"has_scheduled_tasks" form:"has_scheduled_tasks"`
}

func (*Deploy) TableName() string {
	return "deploy"
}
