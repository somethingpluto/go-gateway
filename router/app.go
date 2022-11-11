package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/controller/app"
	"go_gateway/global"
	"go_gateway/middleware"
)

func InitAppRouter(router *gin.Engine) {
	appRouter := router.Group("/app")
	appRouter.Use(sessions.Sessions("mySession", global.SessionRedisStore), middleware.RecoveryMiddleware(), middleware.RequestLog(), middleware.SessionAuthMiddleware(), middleware.TranslationMiddleware())
	{
		app.ServiceRegister(appRouter)
	}
}
