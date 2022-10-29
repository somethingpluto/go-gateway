package dto

import (
	"github.com/gin-gonic/gin"
	"go_gateway/public"
)

type ServiceListInput struct {
	Info    string `json:"info" form:"form" comment:"关键词" example:"http" validate:""`
	PageNo  string `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`
	PageNum string `json:"page_num" form:"page_num" comment:"每页条数" example:"10" validate:"required"`
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListOutput struct {
}
