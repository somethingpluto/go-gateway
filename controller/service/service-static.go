package service

import (
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"time"
)

func (service *ServiceController) ServiceStatic(c *gin.Context) {
	params := &dto.ServiceStaticInput{}
	out := &dto.ServiceStaticOutput{}
	err := params.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 获取服务基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 获取服务详情
	//TODO:暂时置空
	_, err = serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	out.Title = serviceInfo.ServiceName

	todayList := []int{}
	for i := 0; i <= time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	out.Today = todayList
	yesterdayList := []int{}
	for i := 0; i <= 23; i++ {
		yesterdayList = append(yesterdayList, 0)
	}
	out.Yesterday = yesterdayList
	middleware.ResponseSuccess(c, out)
}
