package service

import (
	"JenArgo/model/bo"
	"JenArgo/model/vo"
	"JenArgo/settings"
)

var App app

type app struct{}

// Add 新增
func (*app) Add(apps *vo.AppAddRequest) error {
	// 检查 ProjectName（项目）是否存在
	if err := GitLab.CheckProject(apps.GroupId, apps.ProjectName); err != nil {
		return err
	}

	projectInfo := &bo.GitLabRequest{
		GroupName:   apps.GroupName,
		ProjectName: apps.ProjectName,
		Visibility:  apps.Visibility,
		GroupId:     apps.GroupId,
		Description: apps.Description,
	}

	// 检查 Visibility
	if err := GitLab.CheckVisibility(projectInfo.Visibility); err != nil {
		return err
	}

	// 创建 ProjectName 项目
	if err := GitLab.CreateProject(projectInfo); err != nil {
		return err
	}

	// 创建 Jenkins 流水线任务
	if apps.HasJenkins {
		// 存在则创建 Jenkins 流水线
		newPipelineInfo := &bo.Jenkins{
			GroupName:   apps.GroupName,
			JobName:     apps.ProjectName,
			CopyJobName: settings.Conf.CiCd.CopyJobName,
		}

		if err := Jenkins.CreatePipeline(newPipelineInfo); err != nil {
			return err
		}
	}

	return nil
}
