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
	// 添加更新http
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)

	// 添加更新tcp
	group.POST("/service_add_tcp", service.ServiceAddTCP)
	group.POST("/service_update_tcp", service.ServiceUpdateTCP)
	// 添加更新grpc
	group.POST("/service_add_grpc", service.ServiceAddGRPC)
	group.POST("/service_update_grpc", service.ServiceUpdateGRPC)
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
	// 验证ipList 和 weightList 数量是否一直
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WhiteList, "\n")) {
		middleware.ResponseError(c, 2003, errors.New("ipList数量与权重列表数量不相同"))
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx = tx.Begin()

	// 查找同名服务是否存在
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil { // 存在
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	}

	// 查找同样接入前缀是否存在
	httpRule := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	httpRule, err = httpRule.Find(c, tx, httpRule)
	if err == nil { // 存在
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("接入前缀或域名规则已存在"))
		return
	}
	out.HTTPRule = httpRule

	// serviceInfo保存
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

	// httpRule保存
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
		middleware.ResponseError(c, 2005, err)
		return
	}
	out.HTTPRule = httpRuleModel

	tx.Commit()
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

	out := &dao.ServiceDetail{}
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

	// 查找服务是否存在
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil { // 不存在
		tx.Rollback()
		middleware.ResponseError(c, 2000, errors.New("服务不存在"))
		return
	}

	// 获取服务详细信息
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
	out.Info = info

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
	out.HTTPRule = httpRule

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
	out.AccessControl = accessControl

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
	out.LoadBalance = loadBalance

	tx.Commit()
	middleware.ResponseSuccess(c, out)
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
	// 判断ip与权重数量是否一致
	if len(strings.Split(params.IpList, ",")) == len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
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
	if err == nil { // 占用
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

// ServiceUpdateTCP godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (service *ServiceController) ServiceUpdateTCP(c *gin.Context) {
	out := &dao.ServiceDetail{}
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	tx := lib.GORMDefaultPool.Begin()

	// 获取服务详细信息
	serviceInfo := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := serviceInfo.ServiceDetail(c, lib.GORMDefaultPool, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("服务不存在"))
		return
	}
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}
	out.Info = detail.Info

	// 更行loadBalance
	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}
	out.LoadBalance = loadBalance

	// 更新tcp规则
	tcpRule := &dao.TcpRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	out.TCPRule = tcpRule

	// 更新accessControl
	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	out.AccessControl = accessControl

	tx.Commit()
	middleware.ResponseSuccess(c, out)
}

// ServiceAddGRPC godoc
// @Summary grpc服务添加
// @Description grpc服务添加
// @Tags 服务管理
// @ID /service/service_add_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_grpc [post]
func (service *ServiceController) ServiceAddGRPC(c *gin.Context) {
	out := &dao.ServiceDetail{}
	params := &dto.ServiceAddGrpcInput{}
	err := params.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 验证serviceName是否被占用
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	_, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用"))
		return
	}

	// 验证端口是否被占用
	tcpRule := &dao.TcpRule{Port: params.Port}
	_, err = tcpRule.Find(c, tx, tcpRule)
	if err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端被占用"))
		return
	}
	grpcRule := &dao.GrpcRule{
		Port: params.Port,
	}
	_, err = grpcRule.Find(c, tx, grpcRule)
	if err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端被占用"))
		return
	}

	// 添加serviceInfo信息
	tx = tx.Begin()
	serviceInfo = &dao.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
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

	// 添加loadBalance信息
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

	// 添加grpc规则
	grpcRule = &dao.GrpcRule{
		ServiceID:      serviceInfo.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	err = grpcRule.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	out.GRPCRule = grpcRule
	// accessControl
	accessControl := &dao.AccessControl{
		ServiceID:         serviceInfo.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	out.AccessControl = accessControl

	tx.Commit()
	middleware.ResponseSuccess(c, out)
}

// ServiceUpdateGRPC godoc
// @Summary grpc服务更新
// @Description grpc服务更新
// @Tags 服务管理
// @ID /service/service_update_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_grpc [post]
func (service *ServiceController) ServiceUpdateGRPC(c *gin.Context) {
	out := &dao.ServiceDetail{}
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	tx := lib.GORMDefaultPool.Begin()

	serviceInfo := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := serviceInfo.ServiceDetail(c, lib.GORMDefaultPool, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}
	out.Info = info
	// 负载均衡更新保存
	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	out.LoadBalance = loadBalance

	// grpc规则保存
	grpcRule := &dao.GrpcRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	//grpcRule.Port = params.Port
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	out.GRPCRule = grpcRule

	// accessControl更新保存
	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	tx.Commit()
	out.AccessControl = accessControl
	middleware.ResponseSuccess(c, out)
}
