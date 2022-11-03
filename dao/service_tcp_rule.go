package dao

import (
	"go_gateway/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port      int   `json:"port" gorm:"column:port" description:"端口	"`
}

func (t *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *TcpRule) Find(c *gin.Context, tx *gorm.DB, search *TcpRule) (*TcpRule, error) {
	model := &TcpRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(&model).Error
	return model, err
}
