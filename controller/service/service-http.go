package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"strings"
)

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
