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
	grpcRule.Port = params.Port
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
