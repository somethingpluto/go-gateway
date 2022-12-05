package dashboard

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
	"time"
)

type DashboardController struct {
}

func ServiceRegister(router *gin.RouterGroup) {
	service := &DashboardController{}
	router.GET("/panel_group_data", service.PanelGroupData)
	router.GET("/flow_stat_data", service.FlowStatData)
	router.GET("/pie_chart_data", service.PieChartData)
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

func (service *DashboardController) FlowStatData(c *gin.Context) {
	out := &dto.ServiceStaticOutput{}
	out.Title = "今日流量统计"
	todayList := []int{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 1)
	}
	out.Today = todayList
	yesterdayList := []int{}
	for i := 0; i < 23; i++ {
		yesterdayList = append(yesterdayList, 2)
	}
	out.Yesterday = yesterdayList
	middleware.ResponseSuccess(c, out)
}

func (service *DashboardController) PieChartData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2003, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}

	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}
