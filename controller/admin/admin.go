package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	admin := AdminController{}
	group.GET("/admin_info", admin.AdminInfo)
	group.POST("/change_pwd", admin.ChangePassword)
}

func (admin *AdminController) AdminInfo(c *gin.Context) {
	// 1.读取sessionKey对应的结构体
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	err := json.Unmarshal([]byte(fmt.Sprintf("%s", sessionInfo)), adminSessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
	}
	// 2.取出数据封装输出结构体
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		UserName:     adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "www.baidu.com",
		Introduction: "我是管理员",
		Roles:        []string{"111", "222"},
	}
	middleware.ResponseSuccess(c, out)
}

func (admin AdminController) ChangePassword(c *gin.Context) {
	// 1.获取参数 修改后的密码
	params := &dto.ChangePasswordInput{}
	err := public.DefaultGetValidParams(c, params)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// 2. 通过session获得用户的信息
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	err = json.Unmarshal([]byte(fmt.Sprintf("%s", sessionInfo)), adminSessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 3.查询用户是否存在
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	// 4.重新含盐的密码
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.NewPassword)
	adminInfo.Password = saltPassword
	err = adminInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}
	middleware.ResponseSuccess(c, "密码修改成功")
}
