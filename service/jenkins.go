package service

import (
	"JenArgo/middleware"
	"JenArgo/model/bo"
	"JenArgo/settings"
	"fmt"
	"strings"
)

var Jenkins jenkins

type jenkins struct{}

func (*jenkins) CreatePipeline(data *bo.Jenkins) error {
	// 判断 cocos 安卓 ios 项目和前后端不在一个 jenkins 上
	if strings.Contains(data.JobName, "cocos") || strings.Contains(data.JobName, "android") || strings.Contains(data.JobName, "ios") {
		url := fmt.Sprintf("%s/createItem?name=%s&mode=copy&from=%s", settings.Conf.CiCd.CocosJenkinsUrl, fmt.Sprintf("test-%s-%s", data.GroupName, data.JobName), data.CopyJobName)

		_, err := middleware.Request.HttpRequest("POST", "jenkins", url, settings.Conf.CiCd.CocosUserPassword, nil)
		if err != nil {
			return err
		}
	} else {
		testUrl := fmt.Sprintf("https://test-%s", settings.Conf.CiCd.JenkinsUrl)
		testUrl = fmt.Sprintf("%s/createItem?name=%s&mode=copy&from=%s", testUrl, fmt.Sprintf("%s-%s", data.GroupName, data.JobName), data.CopyJobName)

		_, err := middleware.Request.HttpRequest("POST", "jenkins", testUrl, settings.Conf.CiCd.UserPassword, nil)
		if err != nil {
			return err
		}

		prodUrl := fmt.Sprintf("https://prod-%s", settings.Conf.CiCd.JenkinsUrl)
		prodUrl = fmt.Sprintf("%s/createItem?name=%s&mode=copy&from=%s", prodUrl, fmt.Sprintf("%s-%s", data.GroupName, data.JobName), data.CopyJobName)

		_, err = middleware.Request.HttpRequest("POST", "jenkins", prodUrl, settings.Conf.CiCd.ProdUserPassword, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
