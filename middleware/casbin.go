package middleware

import (
	"JenArgo/common"
	"JenArgo/model/vo"
	"JenArgo/settings"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RbacAuth 权限校验
func RbacAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取用户 Token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			zap.L().Warn("请求未携带 Authorization 头")
			common.ResponseRbacInvalid(ctx, "请求未携带 Authorization 头")
			ctx.Abort()
			return
		}

		// 构造 RBAC 校验请求
		reqBody := &vo.CasbinAuthRequest{
			Path:   ctx.Request.URL.Path,
			Method: ctx.Request.Method,
		}

		// 调用统一的 HTTP 请求方法
		url := fmt.Sprintf("%s/api/rbac/auth", settings.Conf.Rbac.RbacUrl)
		body, err := Request.HttpRequest("POST", "rbac", url, authHeader, reqBody)
		if err != nil {
			zap.L().Error("RBAC HTTP 请求失败", zap.Error(err))
			common.ResponseRbacInvalid(ctx, "rbac service error")
			ctx.Abort()
			return
		}

		// 解析 RBAC 服务响应
		var resp vo.CasbinAuthResponse

		if err := json.Unmarshal(body, &resp); err != nil {
			zap.L().Error("解析 RBAC 响应失败", zap.Error(err))
			common.ResponseRbacInvalid(ctx, "invalid rbac response")
			ctx.Abort()
			return
		}

		// 根据响应码判断权限
		if resp.Code != 0 {
			common.ResponseFailed(ctx, resp.Msg)
			ctx.Abort()
			return
		}

		// 校验通过，继续处理
		ctx.Next()
	}
}
