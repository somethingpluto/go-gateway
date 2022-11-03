package controller

import (
	"fmt"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"

	"github.com/gin-gonic/gin"
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
	out.Total = int(total)
	for _, item := range list {
		serviceDetail, err := item.ServiceDetail(c, db, &item)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := "unknow"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}
		temp := &dto.ServiceListItemOutput{
			ID:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			LoadType:    item.LoadType,
			ServiceAddr: serviceAddr,
		}
		out.List = append(out.List, temp)
	}
	middleware.ResponseSuccess(c, out)
}
