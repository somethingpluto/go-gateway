package dao

import (
	"errors"
	"github.com/e421083458/gorm"

	"github.com/gin-gonic/gin"
	"go_gateway/dto"
	"go_gateway/public"
	"time"
)

type Admin struct {
	Id       int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt     string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password string    `json:"password" gorm:"column:password" description:"密码"`
	UpdateAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (admin *Admin) TableName() string {
	return "gateway_admin"
}

func (admin *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := admin.Find(c, tx, &Admin{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	saltPassword := public.GenSaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码输入错误,请重新输入")
	}
	return adminInfo, nil
}

func (admin *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	err := tx.SetCtx(public.GetTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (admin *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetTraceContext(c)).Save(admin).Error
}
