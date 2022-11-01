package controller

import (
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"strconv"
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
	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, db, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	out := &dto.ServiceListOutput{}
	out.Total = strconv.FormatInt(total, 10)
	for _, item := range list {
		temp := &dto.ServiceListItemOutput{
			ID:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			LoadType:    item.LoadType,
		}
		out.List = append(out.List, temp)
	}
	middleware.ResponseSuccess(c, out)
}
