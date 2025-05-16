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

	//r.Use(middleware.RbacAuth())

	//gitlab
	gitlab := r.Group("/api/gitlab/")
	gitlab.GET("/groups", middleware.RbacAuth(), GitLab.GetGroupsList)
	gitlab.GET("/projects", middleware.RbacAuth(), GitLab.GetProjectsList)
	gitlab.GET("/project/branch", middleware.RbacAuth(), GitLab.GetProjectBranchList)

	//app
	app := r.Group("/api/app/")
	app.POST("/add", App.Add)

	//deploy
	deploy := r.Group("/api/deploy/")
	deploy.GET("/list", middleware.RbacAuth(), Deploy.List)
	deploy.POST("/update", middleware.RbacAuth(), Deploy.Update)
	deploy.POST("/del", middleware.RbacAuth(), Deploy.Delete)
	deploy.POST("/add", Deploy.Add)

	//cicd
	cicd := r.Group("/api/cicd/")
	cicd.POST("/deployCiCd", middleware.RbacAuth(), CiCd.DeployCiCd)
	cicd.POST("/jenkinsCiCd", CiCd.JenkinsCiCd)
	cicd.POST("/updateCiCd", CiCd.UpdateCiCd)

	//argocd
	argocd := r.Group("/api/argocd/")
	argocd.POST("/session", middleware.RbacAuth(), ArgoCD.Session)
	argocd.GET("/apps", middleware.RbacAuth(), ArgoCD.Apps)
	argocd.GET("/image", middleware.RbacAuth(), ArgoCD.Image)
	argocd.POST("/rollback", middleware.RbacAuth(), ArgoCD.Rollback)
	argocd.GET("/log", middleware.RbacAuth(), ArgoCD.Log)

	return r
}
