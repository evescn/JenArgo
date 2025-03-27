package service

import (
	"JenArgo/middleware"
	"JenArgo/model/bo"
	"JenArgo/model/po"
	"JenArgo/model/vo"
	"JenArgo/settings"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"regexp"
)

var GitLab gitlab

type gitlab struct{}

// CheckProject 检查项目是否已存在
func (*gitlab) CheckProject(groupId uint, projectName string) error {
	// 检查 projectName 是否符合规则
	if matched, _ := regexp.MatchString("^[a-z0-9-]+$", projectName); !matched {
		msg := "Error: ProjectName 只能使用小写字母、数字或-\n"
		zap.L().Error(msg)
		return errors.New(msg)
	}

	// 检查 projectName 是否存在
	url := fmt.Sprintf("%s/api/v4/groups/%d/projects?search=%s", settings.Conf.GitLab.GitLabUrl, groupId, projectName)
	body, err := middleware.Request.HttpRequest("GET", "gitlab", url, settings.Conf.GitLab.GitLabToken, nil)
	if err != nil {
		return err
	}

	projectInfo := bo.ProjectInfo{}
	if err := json.Unmarshal(body, &projectInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return err
	}

	if len(projectInfo) != 0 {
		msg := fmt.Sprintf("ERROR: %v 已存在此项目!\n", projectName)
		zap.L().Error(msg)
		return errors.New(msg)
	}

	return nil
}

// CheckVisibility 仓库权限校验
func (*gitlab) CheckVisibility(visibility string) error {
	if visibility != "private" && visibility != "internal" && visibility != "public" {
		msg := "ERROR: 可见度级别关键字错误，只接收 private、internal、public 级别"
		zap.L().Error(msg)
		return errors.New(msg)
	}
	return nil
}

// CreateProject 创建项目
func (*gitlab) CreateProject(args *bo.GitLabRequest) error {
	newProjectInfo := &po.GitLabProjectInfo{
		Name:        args.ProjectName,
		Description: args.Description,
		Path:        args.ProjectName,
		NameSpaceId: args.GroupId,
		Visibility:  args.Visibility,
		ImportUrl:   settings.Conf.GitLab.GitLabUrl + "/init/bare.git",
	}

	url := settings.Conf.GitLab.GitLabUrl + "/api/v4/projects"
	body, err := middleware.Request.HttpRequest("POST", "gitlab", url, settings.Conf.GitLab.GitLabToken, newProjectInfo)
	if err != nil {
		return err
	}

	zap.L().Info("创建项目成功: " + string(body))
	return nil
}

// GetGroupsList 获取组信息
func (*gitlab) GetGroupsList() (id interface{}, err error) {
	url := settings.Conf.GitLab.GitLabUrl + "/api/v4/groups"
	body, err := middleware.Request.HttpRequest("GET", "gitlab", url, settings.Conf.GitLab.GitLabToken, nil)
	if err != nil {
		return 0, err
	}

	groupInfo := bo.GroupsInfo{}
	if err := json.Unmarshal(body, &groupInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return 0, err
	}

	return groupInfo, nil
}

// GetProjectsList 获取项目信息
func (*gitlab) GetProjectsList(projectInfo *vo.GitLabProjectListRequest) (*bo.ProjectInfo, error) {
	var url string
	if projectInfo.GroupId != 0 {
		url = fmt.Sprintf("%s/api/v4/groups/%d/search?scope=projects&search=%s&page=%d&per_page=%d", settings.Conf.GitLab.GitLabUrl, projectInfo.GroupId, projectInfo.AppName, projectInfo.Page, projectInfo.Size)
	} else {
		url = fmt.Sprintf("%s/api/v4/search?scope=projects&search=%s&page=%d&per_page=%d", settings.Conf.GitLab.GitLabUrl, projectInfo.AppName, projectInfo.Page, projectInfo.Size)
	}

	body, err := middleware.Request.HttpRequest("GET", "gitlab", url, settings.Conf.GitLab.GitLabToken, nil)
	if err != nil {
		return nil, err
	}

	projectsInfo := &bo.ProjectInfo{}
	if err := json.Unmarshal(body, projectsInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return nil, err
	}

	return projectsInfo, nil
}

// GetProjectBranchOrTagList 获取项目分支信息
func (*gitlab) GetProjectBranchOrTagList(appId uint, branchOrTag bool) (*bo.BranchInfo, error) {
	var url string
	if branchOrTag {
		url = fmt.Sprintf("%s/api/v4/projects/%d/repository/tags", settings.Conf.GitLab.GitLabUrl, appId)
	} else {
		url = fmt.Sprintf("%s/api/v4/projects/%d/repository/branches", settings.Conf.GitLab.GitLabUrl, appId)
	}

	body, err := middleware.Request.HttpRequest("GET", "gitlab", url, settings.Conf.GitLab.GitLabToken, nil)
	if err != nil {
		return nil, err
	}

	projectBranchInfo := &bo.BranchInfo{}
	if err := json.Unmarshal(body, projectBranchInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return nil, err
	}

	return projectBranchInfo, nil
}
