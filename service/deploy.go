package service

import (
	"JenArgo/dao"
	"JenArgo/middleware/snowflake"
	"JenArgo/model/vo"
	"JenArgo/settings"
	"errors"
	"fmt"
)

var Deploy deploy

type deploy struct{}

// List 列表
func (*deploy) List(en, appName, repoName string, page, limit int) (*vo.DeploysListResponse, error) {
	return dao.Deploy.List(en, appName, repoName, page, limit)
}

// Add 新增
func (*deploy) Add(data *vo.DeployRequest) error {
	// 拼接下 BuildUrl 信息
	data.ID = snowflake.GenID()
	if data.En == "prod" {
		data.BuildUrl = fmt.Sprintf("https://prod-%s", settings.Conf.CiCd.JenkinsUrl)
	} else {
		data.BuildUrl = fmt.Sprintf("https://test-%s", settings.Conf.CiCd.JenkinsUrl)
	}
	data.BuildUrl = fmt.Sprintf("%s/view/%s/job/%s-%s/", data.BuildUrl, data.RepoName, data.RepoName, data.AppName)

	// 判断是否已存在此数据
	d, has, err := dao.Deploy.Has(data.En, data.AppName, data.RepoName)
	if err != nil {
		return err
	}
	deployInfo := data.ToDeploy()
	if !has {
		return dao.Deploy.Add(deployInfo)
	}

	if d.Status == 1 {
		return errors.New("数据已存在")
	}
	return dao.Deploy.Add(deployInfo)
}

// Update 更新
func (*deploy) Update(d *vo.DeployRequest) error {
	return dao.Deploy.Update(d.ToDeploy())
}

// Delete 删除
func (*deploy) Delete(id int64) error {
	return dao.Deploy.Delete(id)
}
