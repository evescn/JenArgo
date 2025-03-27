package controller

import (
	"JenArgo/common"
	"JenArgo/model/vo"
	"JenArgo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var CiCd cicd

type cicd struct{}

// DeployCiCd 新增
func (*cicd) DeployCiCd(ctx *gin.Context) {
	params := new(vo.DeployRequest)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.CiCd.DeployCiCd(params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "构建部署任务触发")
}

// JenkinsCiCd 新增
func (*cicd) JenkinsCiCd(ctx *gin.Context) {
	params := new(vo.JenArgoCiCd)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.CiCd.JenkinsCiCd(params.En, params.AppName, params.RepoName, params.Builder)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "新增部署任务触发")
}

// UpdateCiCd 列表
func (*cicd) UpdateCiCd(ctx *gin.Context) {
	params := new(vo.JenArgoUpdateCiCd)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("Bind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.CiCd.UpdateCiCd(params.En, params.AppName, params.RepoName, params.Branch, params.CodeCheck, params.BuildStatus, params.DeployStatus)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "更新 CICD 流程成功")
}
