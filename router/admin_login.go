package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/controller/admin"
	"go_gateway/global"
	"go_gateway/middleware"
)

func InitAdminLoginRouter(router *gin.Engine) {
	adminLoginRouter := router.Group("/admin_login")
	adminLoginRouter.Use(sessions.Sessions("mySession", global.SessionRedisStore), middleware.RecoveryMiddleware(), middleware.RequestLog(), middleware.TranslationMiddleware())
	{
		admin.AdminLoginRegister(adminLoginRouter)
	}
}
