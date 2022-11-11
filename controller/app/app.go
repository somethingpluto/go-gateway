package app

import (
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
)

type AppController struct{}

func ServiceRegister(router *gin.RouterGroup) {
	app := &AppController{}
	router.GET("/app_list", app.AppList)
}

// AppList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
func (service *AppController) AppList(c *gin.Context) {
	params := &dto.APPListInput{}
	err := params.GetValidParams(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	appInfo := &dao.App{}
	list, total, err := appInfo.APPList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	outPutList := []dto.APPListItemOutput{}
	for _, item := range list {
		outPutList = append(outPutList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  0,
			RealQps:  0,
		})
	}

	out := &dto.APPListOutput{
		List:  outPutList,
		Total: total,
	}
	middleware.ResponseSuccess(c, out)
}
