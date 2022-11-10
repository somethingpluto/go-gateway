package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/controller/admin"
	"go_gateway/global"
	"go_gateway/middleware"
)

func InitAdminRouter(router *gin.Engine) {
	adminRouter := router.Group("/admin")
	adminRouter.Use(sessions.Sessions("mySession", global.SessionRedisStore), middleware.RecoveryMiddleware(), middleware.RequestLog(), middleware.SessionAuthMiddleware(), middleware.TranslationMiddleware())
	{
		admin.AdminRegister(adminRouter)
	}
}
