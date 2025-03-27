package controller

import (
	"JenArgo/common"
	"JenArgo/model/vo"
	"JenArgo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Deploy deploy

type deploy struct{}

// List 列表
func (*deploy) List(ctx *gin.Context) {
	//接收参数
	params := new(vo.DeploysListRequest)

	//绑定参数
	if err := ctx.Bind(params); err != nil {
		zap.L().Error("Bind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	data, err := service.Deploy.List(params.En, params.AppName, params.RepoName, params.Page, params.Size)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}
	common.ResponseOk(ctx, data, "获取部署列表成功")
}

// Update 更新
func (*deploy) Update(ctx *gin.Context) {
	params := new(vo.DeployRequest)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.Deploy.Update(params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}
	common.ResponseOk(ctx, nil, "更新部署信息成功")
}

// Add 新增
func (*deploy) Add(ctx *gin.Context) {
	params := new(vo.DeployRequest)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.Deploy.Add(params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}
	common.ResponseOk(ctx, nil, "新增部署成功")
}

// Delete 删除
func (*deploy) Delete(ctx *gin.Context) {
	params := new(vo.DeployDelRequest)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.Deploy.Delete(params.ID)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}

	common.ResponseOk(ctx, nil, "删除部署信息成功")
}
