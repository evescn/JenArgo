package service

import (
	"JenArgo/dao"
	"JenArgo/middleware"
	"JenArgo/model/vo"
	"JenArgo/settings"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

var CiCd cicd

type cicd struct{}

// 使用 sync.WaitGroup 防止主进程退出，定时任务被回收
var wg sync.WaitGroup

// DeployCiCd 开始部署
func (c *cicd) DeployCiCd(d *vo.DeployRequest) error {
	data, has, err := dao.Deploy.Get(d.ID)
	if err != nil {
		return err
	}

	if !has {
		return errors.New("查询无此部署任务")
	}

	// 更新数据
	d.Status = 1
	if data.HasScheduledTask && data.StartTime != "" {
		// 启用了定时器 解析时间字符串
		targetTime, err := time.ParseInLocation("2006-01-02 15:04:05", data.StartTime, time.Local)
		if err != nil {
			return err
		}

		// 计算当前时间到目标时间的间隔
		delay := time.Until(targetTime)
		if delay <= 0 {
			return errors.New("目标时间已过，不执行任务")
		}

		err = dao.Deploy.Update(d.ToDeploy())
		if err != nil {
			return err
		}

		// 设置定时器 请求 jenkins 服务
		wg.Add(1)
		time.AfterFunc(delay, func() {
			defer wg.Done()
			zap.L().Info(fmt.Sprintf("定时任务 %d 触发，正在请求 Jenkins...", d.ID))

			// 请求 jenkins 服务
			if err := c.triggerJenkins(data); err != nil {
				zap.L().Error(fmt.Sprintf("任务 %d 触发失败: %v", d.ID, err.Error()))
			} else {
				log.Printf("任务 %d 触发成功", d.ID)
				zap.L().Info(fmt.Sprintf("任务 %d 触发成功", d.ID))
			}
		})

	} else {
		// 没有启用定时器 获取当前本地时间 转换为字符串格式
		now := time.Now().Local()
		d.StartTime = now.Format("2006-01-02 15:04:05")

		err = dao.Deploy.Update(d.ToDeploy())
		if err != nil {
			return err
		}

		// 请求 jenkins 服务
		if err := c.triggerJenkins(data); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (*cicd) triggerJenkins(task *vo.DeployRequest) error {
	url := fmt.Sprintf("%sbuildWithParameters?ENV=%s&BRANCH=%s&IS_CREATE_TAG=%v",
		task.BuildUrl, task.En, task.Branch, task.Tag)
	if task.En == "prod" {
		_, err := middleware.Request.HttpRequest("POST", "jenkins", url, settings.Conf.CiCd.ProdUserPassword, nil)
		if err != nil {
			return err
		}
	} else {
		_, err := middleware.Request.HttpRequest("POST", "jenkins", url, settings.Conf.CiCd.UserPassword, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// JenkinsCiCd 开始部署
func (*cicd) JenkinsCiCd(en, appName, repoName, builder string) error {
	data, has, err := dao.Deploy.Has(en, appName, repoName)
	if err != nil {
		return err
	}

	if !has {
		return errors.New("查询无此部署任务")
	}
	if data.Status == 1 {
		return errors.New("此任务正在部署中")
	}

	// 更新数据
	data.Status = 1
	now := time.Now().Local()
	data.StartTime = now.Format("2006-01-02 15:04:05")
	data.Builder = builder

	err = dao.Deploy.Update(data)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCiCd 更新CiCd流程
func (cicd *cicd) UpdateCiCd(en, appName, repoName, branch string, codeCheck, buildStatus, deployStatus int) error {
	var deployID int64

	// 使用 en 和 appName 获取数据
	deploysData, err := dao.Deploy.List(en, appName, repoName, 1, 10)
	if err != nil {
		return err
	}

	// 返回第一条记录即可
	for _, item := range deploysData.Items {
		deployID = item.ID
		break
	}

	// 获取数据
	deployData, _, err := dao.Deploy.Get(deployID)
	if err != nil {
		return err
	}

	if deployID == 0 {
		return errors.New("查询无此应用")
	}

	// 判断 codeCheck 是否 != 0
	if codeCheck != 0 {
		deployData.CodeCheck = codeCheck
		if codeCheck == 2 {
			deployData.Status = 4
			// 部署耗时
			deployData.Duration = cicd.getDuration(deployData.StartTime)
		}
	}

	// 判断 buildStatus 是否 == 1
	if buildStatus == 1 {
		deployData.BuildStatus = buildStatus
		deployData.Status = 2
	} else if buildStatus == 2 {
		deployData.BuildStatus = buildStatus
		deployData.Status = 4
		// 部署耗时
		deployData.Duration = cicd.getDuration(deployData.StartTime)

	}

	// 判断 deployStatus 是否 == 2
	if deployStatus == 1 && deployData.Status == 2 {
		deployData.DeployStatus = deployStatus
		deployData.Status = 3
		// 程序启动日志查看
		// 部署耗时
		deployData.Duration = cicd.getDuration(deployData.StartTime)
	} else if deployStatus == 2 {
		deployData.DeployStatus = deployStatus
		deployData.Status = 4
		// 部署耗时
		deployData.Duration = cicd.getDuration(deployData.StartTime)
	}

	// 更新 Tag 信息
	if len(branch) > 0 {
		deployData.Branch = branch
	}

	// 更新 Tag 信息
	if len(branch) > 0 {
		deployData.Branch = branch
	}

	// 更新服务
	return dao.Deploy.Update(deployData.ToDeploy())
}

// getDuration 获取时间，计算耗时
func (*cicd) getDuration(startTimeStr string) string {
	// 解析开始时间
	loc, _ := time.LoadLocation("Asia/Shanghai")

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, loc)
	if err != nil {
		zap.L().Error(err.Error())
		return ""
	}

	// 获取当前时间作为结束时间
	endTime := time.Now()

	// 计算时间差
	duration := endTime.Sub(startTime)
	// 使用 String 方法将时间差格式化为分秒形式
	durationString := duration.Truncate(time.Second).String()

	// 打印时间差
	return durationString
}
