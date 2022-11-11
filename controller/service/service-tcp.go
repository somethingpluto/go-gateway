package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
	"strings"
)

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
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
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
	grpcRule := &dao.GrpcRule{Port: params.Port}
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
