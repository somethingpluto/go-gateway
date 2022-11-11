package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/controller/dashboard"
	"go_gateway/global"
	"go_gateway/middleware"
)

func InitDashboardRouter(router *gin.Engine) {
	dashboardRouter := router.Group("/dashboard")
	dashboardRouter.Use(sessions.Sessions("mySession", global.SessionRedisStore), middleware.RecoveryMiddleware(), middleware.SessionAuthMiddleware(), middleware.TranslationMiddleware())
	{
		dashboard.ServiceRegister(dashboardRouter)
	}
}
