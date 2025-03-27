package controller

import (
	"JenArgo/common"
	"JenArgo/model/vo"
	"JenArgo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var GitLab gitlab

type gitlab struct{}

func (*gitlab) GetGroupsList(ctx *gin.Context) {
	data, err := service.GitLab.GetGroupsList()
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "获取 Groups 成功")
}

func (*gitlab) GetProjectsList(ctx *gin.Context) {
	param := &vo.GitLabProjectListRequest{}

	if err := ctx.Bind(param); err != nil {
		zap.L().Error("Bind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}
	data, err := service.GitLab.GetProjectsList(param)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "获取 Projects 成功")
}

func (*gitlab) GetProjectBranchList(ctx *gin.Context) {
	param := &vo.GitLabProjectBranchRequest{}

	if err := ctx.Bind(param); err != nil {
		zap.L().Error("Bind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.GitLab.GetProjectBranchOrTagList(param.AppId, param.BranchOrTag)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "获取 Project 分支成功")
}
