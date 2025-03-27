package controller

import (
	"JenArgo/logger"
	"JenArgo/middleware"
	"JenArgo/settings"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router router

type router struct{}

func (*router) Setup() *gin.Engine {
	// 初始化gin对象
	if settings.Conf.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// 修改日志格式
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 跨域中间件
	r.Use(middleware.Cross.Cors())
	// JWT登陆验证中间件
	//r.Use(middle.JWTAuth())

	r.GET("/testApi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "testApi success!",
			"data": nil,
		})
	})

	r.GET("/version", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, settings.Conf.Version)
	})

	//gitlab
	gitlab := r.Group("/api/gitlab/")
	gitlab.GET("/groups", GitLab.GetGroupsList)
	gitlab.GET("/projects", GitLab.GetProjectsList)
	gitlab.GET("/project/branch", GitLab.GetProjectBranchList)

	//app
	app := r.Group("/api/app/")
	app.POST("/add", App.Add)

	//deploy
	deploy := r.Group("/api/deploy/")
	deploy.GET("/list", Deploy.List)
	deploy.POST("/update", Deploy.Update)
	deploy.POST("/del", Deploy.Delete)
	deploy.POST("/add", Deploy.Add)

	//cicd
	cicd := r.Group("/api/cicd/")
	cicd.POST("/deployCiCd", CiCd.DeployCiCd)
	cicd.POST("/jenkinsCiCd", CiCd.JenkinsCiCd)
	cicd.POST("/updateCiCd", CiCd.UpdateCiCd)

	//argocd
	argocd := r.Group("/api/argocd/")
	argocd.POST("/session", ArgoCD.Session)
	argocd.GET("/apps", ArgoCD.Apps)
	argocd.GET("/image", ArgoCD.Image)
	argocd.POST("/rollback", ArgoCD.Rollback)
	argocd.GET("/log", ArgoCD.Log)

	return r
}
