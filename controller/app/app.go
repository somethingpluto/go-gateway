package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
)

type APPController struct{}

func ServiceRegister(router *gin.RouterGroup) {
	app := &APPController{}
	router.GET("/app_list", app.APPList)
	router.GET("/app_detail", app.APPDetail)
	router.DELETE("/app_delete", app.APPDelete)
	router.POST("/app_add", app.APPAdd)
}

// APPList godoc
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
// @Router /app/app_list [get]
func (service *APPController) APPList(c *gin.Context) {
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
	var outPutList []dto.APPListItemOutput
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

// APPDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (service *APPController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	err := params.GetValidParams(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	search := &dao.App{ID: params.ID}
	detail, err := search.Find(c, tx, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, detail)
}

// APPDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (service *APPController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	err := params.GetValidParams(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{ID: params.ID}
	appInfo, err := search.Find(c, tx, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	appInfo.IsDelete = 1
	err = appInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

// APPAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (service *APPController) APPAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	err := params.GetValidParams(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 判断id 是否被占用
	appInfo := &dao.App{
		AppID: params.AppID,
	}
	_, err = appInfo.Find(c, tx, appInfo)
	if err == nil {
		middleware.ResponseError(c, 2002, errors.New("租户ID已被占用"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	appInfo = &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	err = appInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [post]
func (service *APPController) AppUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	err := params.GetValidParams(c)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	appInfo, err := search.Find(c, tx, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	appInfo.Name = params.Name
	appInfo.Secret = params.Secret
	appInfo.WhiteIPS = params.WhiteIPS
	appInfo.Qps = params.Qps
	appInfo.Qpd = params.Qpd
	err = appInfo.Save(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}
