package dto

import (
	"github.com/gin-gonic/gin"
	"go_gateway/public"
)

type ServiceStaticInput struct {
	ID int64 `json:"id" form:"id" validate:"required" comment:"服务id"`
}

type ServiceStaticOutput struct {
	Today     []int `json:"today"`
	Yesterday []int `json:"yesterday"`
}

func (params *ServiceStaticInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
