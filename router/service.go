package router

import (
	"github.com/gin-gonic/gin"
	"go_gateway/controller"
)

func InitServiceRouter(router *gin.Engine) {
	serviceRouter := router.Group("/service")
	{
		controller.ServiceRegister(serviceRouter)
	}
}
