package controller

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
	"time"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)
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
	// 获取用户输入
	param := &dto.AdminLoginInput{}
	// 用户输入验证
	err := param.BindValidParam(c)
	if err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	// 1. 获取用户登录信息
	//userName := param.UserName
	//password := param.Password
	// 2. salt + password sha256加密 =>saltPassword
	// 获取db
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	admin := &dao.Admin{}
	// 检查用户密码是否正确
	admin, err = admin.LoginCheck(c, tx, param)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 设置session

	sessionInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessionInfoBytes, err := json.Marshal(sessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	session := sessions.Default(c)
	session.Set(public.AdminSessionInfoKey, string(sessionInfoBytes))
	session.Save()
	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}

func (adminLogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(public.AdminSessionInfoKey)
	session.Save()
	middleware.ResponseSuccess(c, "退出登录")
}
