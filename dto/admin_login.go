package dto

import (
	"github.com/gin-gonic/gin"
	"go_gateway/public"
	"time"
)

// AdminLoginInput
// @Description: 管理员登录输入
//
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"姓名" example:"admin" validate:"required"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	err := public.DefaultGetValidParams(c, param)
	return err
}

// AdminLoginOutput
// @Description: 管理登录 输出
//
type AdminLoginOutput struct {
	Token string `json:"token"`
}

type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}
