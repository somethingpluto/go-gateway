package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/controller"
	"go_gateway/middleware"
	"log"
)

func InitAdminLoginRouter(router *gin.Engine) {
	adminLoginRouter := router.Group("/admin_login")
	store, err := sessions.NewRedisStore(10, "tcp", "120.25.255.207:6380", "chx200205173214", []byte("secret"))
	if err != nil {
		log.Fatalln(err)
	}

	adminLoginRouter.Use(sessions.Sessions("mySession", store), middleware.RecoveryMiddleware(), middleware.RequestLog(), middleware.TranslationMiddleware())
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}
}
