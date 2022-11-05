package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
	"strings"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.POST("/service_delete", service.ServiceDelete)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
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
		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		temp := &dto.ServiceListItemOutput{
			ID:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			LoadType:    item.LoadType,
			ServiceAddr: serviceAddr,
			Qpd:         0,
			Qps:         0,
			TotalNode:   len(ipList),
		}
		out.List = append(out.List, temp)
	}
	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	out := &dao.ServiceDetail{}
	params := &dto.ServiceDeleteInput{}
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

	// 读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	serviceInfo.IsDelete = 1
	err = serviceInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	out.Info.ID = serviceInfo.ID
	middleware.ResponseSuccess(c, out)
}

// ServiceAddHTTP godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (service *ServiceController) ServiceAddHTTP(c *gin.Context) {
	out := &dao.ServiceDetail{}
	params := &dto.ServiceAddHTTPInput{}
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

	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	// 查找同名服务是否存在
	tx.Begin()
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	}

	httpRule := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	httpRule, err = httpRule.Find(c, tx, httpRule)
	if err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("接入前缀或域名规则已存在"))
		return
	}
	out.HTTPRule = httpRule
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WhiteList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("ipList数量与权重列表数量不相同"))
		return
	}

	serviceInfoModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	err = serviceInfoModel.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}
	out.Info = serviceInfoModel
	httpRuleModel := &dao.HttpRule{
		ServiceID:      serviceInfoModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	err = httpRuleModel.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}
	out.HTTPRule = httpRuleModel
	middleware.ResponseSuccess(c, out)
}

// ServiceUpdateHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service_update_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_http [post]
func (service *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	err := params.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	tx = tx.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2000, errors.New("服务不存在"))
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2000, errors.New("服务不存在"))
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	err = info.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	// 修改httpRule表
	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	err = httpRule.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 修改accessControl
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	err = accessControl.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	// 修改loadBalance
	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	err = loadBalance.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "更新成功")
}

// ServiceAddTCP godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (service *ServiceController) ServiceAddTCP(c *gin.Context) {
	out := &dao.ServiceDetail{}
	params := &dto.ServiceAddTcpInput{}
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
	// 查询ServiceName是否被占用
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名已被占用"))
		return
	}

	// 验证端口是否被占用
	tcpRule := &dao.TcpRule{Port: params.Port}
	_, err = tcpRule.Find(c, tx, tcpRule)
	if err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口已被占用"))
		return
	}
	grpcRule := &dao.GrpcRule{}
	_, err = grpcRule.Find(c, tx, grpcRule)
	if err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口已被占用"))
		return
	}
	// 判断ip与权重数量是否一致
	if len(strings.Split(params.IpList, ",")) == len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}
	tx = tx.Begin()
	// serviceInfo保存
	serviceInfo = &dao.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	err = serviceInfo.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	out.Info = serviceInfo

	// loadBalance保存
	loadBalance := &dao.LoadBalance{
		ServiceID:  serviceInfo.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	err = loadBalance.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	out.LoadBalance = loadBalance

	// tcp规则保存
	tcpRule = &dao.TcpRule{
		ServiceID: serviceInfo.ID,
		Port:      params.Port,
	}
	err = tcpRule.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	out.TCPRule = tcpRule

	// accessControl保存
	accessControl := &dao.AccessControl{
		ServiceID:         serviceInfo.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	err = accessControl.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	out.AccessControl = accessControl

	// 提交事务
	tx.Commit()
	middleware.ResponseSuccess(c, out)
}
