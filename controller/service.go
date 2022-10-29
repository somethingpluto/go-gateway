package controller

import "github.com/gin-gonic/gin"

type ServiceController struct{}

func ServiceRegister(router *gin.RouterGroup) {
	service := &ServiceController{}
	router.GET("/service_list", service.ServiceList)
}

func (service *ServiceController) ServiceList(c *gin.Context) {

}