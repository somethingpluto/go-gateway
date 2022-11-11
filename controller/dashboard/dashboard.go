package dashboard

import (
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
)

type DashboardController struct {
}

func ServiceRegister(router *gin.RouterGroup) {

}

func (service *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageNo: 1, PageSize: 10})
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          10,
		TodayRequestNum: 10,
		CurrentQPS:      10,
	}
	middleware.ResponseSuccess(c, out)
}
