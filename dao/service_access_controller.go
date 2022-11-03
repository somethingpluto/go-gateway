package dao

import (
	"go_gateway/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// AccessControl
// @Description: gateway_service_access_control 表实体结构体
//
type AccessControl struct {
	ID                int64  `json:"id" gorm:"primary_key"`
	ServiceID         int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip	"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip	"`
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机	"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流	"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端限流	"`
}

func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

// Find
// @Description: 根据条件查找Access
// @receiver t
// @param c
// @param tx
// @param search
// @return *AccessControl
// @return error
//
func (t *AccessControl) Find(c *gin.Context, tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	model := &AccessControl{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(&model).Error
	return model, err
}
