package dto

import (
	"github.com/gin-gonic/gin"
	"go_gateway/public"
)

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"http" validate:""`
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"10" validate:"required"`
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListOutput struct {
	Total int                      `json:"total" form:"total" comment:"总数" example:"" validate:""`
	List  []*ServiceListItemOutput `json:"list" form:"list" comment:"列表" `
}

type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`
	ServiceName string `json:"service_name" form:"service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc"`
	LoadType    int    `json:"load_type" form:"load_type"`
	ServiceAddr string `json:"service_addr" from:"service_addr"`
	Qps         int64  `json:"qps" form:"qps"`
	Qpd         int64  `json:"qpd" form:"qpd"`
	TotalNode   int    `json:"total_node" form:"total_node"`
}

type ServiceDeleteInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"50" validate:"required"`
}

func (param *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddHTTPInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`     //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"0" validate:"max=1,min=0"`                          //接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required,valid_rule"`                      //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"支持https" example:"0" validate:"max=1,min=0"`                     //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"0" validate:"max=1,min=0"`         //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"0" validate:"max=1,min=0"`       //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`              //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换" example:"" validate:"valid_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"0" validate:"max=1,min=0"`                //关键词
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`                           //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`                           //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"0" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"0" validate:"min=0"`     //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"0" validate:"max=3,min=0"`                                //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"0" validate:"required,valid_iplist"`                            //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weightlist"`                 //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时, 单位s" example:"0" validate:"min=0"`   //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时, 单位s" example:"0" validate:"min=0"` //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间, 单位s" example:"0" validate:"min=0"`       //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"10" validate:"min=0"`                    //最大空闲链接数
}

func (param *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateHTTPInput struct {
	ID          int64  `json:"id" form:"id" comment:"服务ID" example:"" validate:"required"`
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`     //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"0" validate:"max=1,min=0"`                          //接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required,valid_rule"`                      //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"支持https" example:"0" validate:"max=1,min=0"`                     //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"0" validate:"max=1,min=0"`         //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"0" validate:"max=1,min=0"`       //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`              //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换" example:"" validate:"valid_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"0" validate:"max=1,min=0"`                //关键词
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`                           //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`                           //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"0" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"0" validate:"min=0"`     //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"0" validate:"max=3,min=0"`                                //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"0" validate:"required,valid_iplist"`                            //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weightlist"`                 //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时, 单位s" example:"0" validate:"min=0"`   //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时, 单位s" example:"0" validate:"min=0"` //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间, 单位s" example:"0" validate:"min=0"`       //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"10" validate:"min=0"`                    //最大空闲链接数
}

func (param *ServiceUpdateHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddTcpInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"header头转换" validate:""`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceAddTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceUpdateTcpInput struct {
	ID                int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceUpdateTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ServiceAddGrpcInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`
	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:"valid_header_transfor"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

func (params *ServiceAddGrpcInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
