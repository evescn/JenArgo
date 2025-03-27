package controller

import (
	"JenArgo/common"
	"JenArgo/model/vo"
	"JenArgo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ArgoCD argoCD

type argoCD struct{}

func (*argoCD) Session(ctx *gin.Context) {
	err := service.ArgoCD.Session()
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "获取 ArgoCD Session 成功")
}

//applications

func (*argoCD) Apps(ctx *gin.Context) {
	params := &vo.ArgoCDApp{}

	if err := ctx.Bind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.ArgoCD.Applications(params.Name, params.Namespace, params.Page, params.Size)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "获取 ArgoCD Apps 成功")
}

func (*argoCD) Image(ctx *gin.Context) {
	params := &vo.ArgoCDImage{}

	if err := ctx.Bind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.ArgoCD.Image(params.Name)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "获取 ArgoCD App 信息成功")
}

func (*argoCD) Rollback(ctx *gin.Context) {
	params := &vo.ArgoCDRollback{}

	if err := ctx.Bind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.ArgoCD.Rollback(params.Name, params.RollbackID)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "App 回滚服务成功")
}

// Log 日志
func (*argoCD) Log(ctx *gin.Context) {
	params := &vo.ArgoCDApp{}

	if err := ctx.Bind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.ArgoCD.Log(params.Name, params.Namespace)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, data, "获取日志成功")
}
