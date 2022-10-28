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

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	param := &dto.AdminLoginInput{}
	err := param.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	// 1. 获取用户登录信息
	//userName := param.UserName
	//password := param.Password
	// 2. salt + password sha256加密 =>saltPassword
	out := &dto.AdminLoginOutput{Token: "token"}
	middleware.ResponseSuccess(c, out)
}
