package controller

import (
	"github.com/e421083458/gin_scaffold/middleware"
	"github.com/gin-gonic/gin"
	"go_gateway/dto"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
}

func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	param := &dto.AdminLoginInput{}
	err := param.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}
