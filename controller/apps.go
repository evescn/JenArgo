package controller

import (
	"JenArgo/common"
	"JenArgo/model/vo"
	"JenArgo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var App app

type app struct{}

// Add 新增
func (*app) Add(ctx *gin.Context) {
	params := new(vo.AppAddRequest)

	if err := ctx.ShouldBind(params); err != nil {
		zap.L().Error("ShouldBind请求参数失败," + err.Error())
		common.ResponseParamInvalid(ctx, err.Error())
		return
	}

	err := service.App.Add(params)
	if err != nil {
		common.ResponseFailed(ctx, err.Error())
		return
	}
	common.ResponseOk(ctx, nil, "新增 GitLab 仓库成功")
}
