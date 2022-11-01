package controller

import (
	"github.com/gin-gonic/gin"
	"go_gateway/dto"
	"go_gateway/middleware"
)

type ServiceController struct{}

func ServiceRegister(router *gin.RouterGroup) {
	service := &ServiceController{}
	router.GET("/service_list", service.ServiceList)
}

func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	err := params.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	out := &dto.ServiceListOutput{}
	middleware.ResponseSuccess(c, out)
}
