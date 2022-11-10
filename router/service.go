package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/controller/service"
	"go_gateway/global"
	"go_gateway/middleware"
)

func InitServiceRouter(router *gin.Engine) {
	serviceRouter := router.Group("/service")
	serviceRouter.Use(sessions.Sessions("mySession", global.SessionRedisStore), middleware.RecoveryMiddleware(), middleware.RequestLog(), middleware.SessionAuthMiddleware(), middleware.TranslationMiddleware())
	{
		service.ServiceRegister(serviceRouter)
	}
}
