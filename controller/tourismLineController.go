package controller

import (
	"github.com/gin-gonic/gin"
	"local/jshw/logic"
	"local/jshw/model/tourismline"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/cihub/seelog"
)

//获取简单的线路列表
func TourismLineList(c *gin.Context) {
	respData := map[string]interface{}{}

	lines := []model.TourismLine{}
	err := logic.GetTourismList(&lines)
	if err != nil {
		log.Println(err)
		ReturnError(c, err)
		return
	}

	respData["suc"] = true
	respData["data"] = lines
	c.JSON(http.StatusOK, respData)
}

//获取线路详情
func TourismLineDetail(c *gin.Context) {
	respData := map[string]interface{}{}
	lineID := c.PostForm("lineID")

	lineIDInt, err := strconv.Atoi(lineID)
	if err != nil {
		log.Println(err)
		ReturnError(c, err)
		return
	}
	tourismLineDetail := model.TourismLineDetail{}

	tl := []model.TourismLine{}
	tld := model.TourismLineDate{TourismLineID:lineIDInt}
	tlt := []model.TourismLineTrip{}
	tli := []model.TourismLineImages{}
	tlp := []model.TourismLinePoint{}
	err = logic.GetTourismLineByID(&tl, lineIDInt)
	if err != nil {
		seelog.Error(err)
		ReturnError(c, err)
		return
	}
	_, err = logic.GetTourismLineDateByLineID(&tld)
	if err != nil {
		seelog.Error(err)
		ReturnError(c, err)
		return
	}
	err = logic.GetTourismLineTripByLineID(&tlt, lineIDInt)
	if err != nil {
		seelog.Error(err)
		ReturnError(c, err)
		return
	}
	err = logic.GetTourismLineImagesByLineID(&tli, lineIDInt)
	if err != nil {
		seelog.Error(err)
		ReturnError(c, err)
		return
	}
	err = logic.GetTourismLinePointByID(&tlp, strings.Split(tl[0].Rendezvous, ","))
	if err != nil {
		seelog.Error(err)
		ReturnError(c, err)
		return
	}

	//格式化路线中的img
	formatTlt := []model.TourismLineTrip{}
	for _, trip := range tlt {
		//replacer:=strings.NewReplacer("src=\"/","<img src=\"http://www.jshwclub.com/")
		replacer := strings.NewReplacer(
			"<img", "",
			"src=\"/upload/", "",
			"width=", "",
			"height=", "",
			"alt=", "",
		)
		tlt := model.TourismLineTrip{}
		tlt.ID = trip.ID
		tlt.Title = trip.Title
		tlt.Contents = replacer.Replace(trip.Contents)
		formatTlt = append(formatTlt, tlt)
	}
	//更改集合点中百度地图iframe的宽度
	replacer := strings.NewReplacer(
		"558", "300", //iframe宽度
		"560", "300", //地图层的宽度
	)
	tl[0].PlaceExplain = replacer.Replace(tl[0].PlaceExplain)

	tourismLineDetail.Line = tl[0]
	tourismLineDetail.LineDate = tld
	tourismLineDetail.LineTrip = formatTlt
	tourismLineDetail.LineImage = tli
	tourismLineDetail.LinePoint = tlp

	respData["suc"] = true
	respData["data"] = tourismLineDetail
	c.JSON(http.StatusOK, respData)
}

//获取线路集合点
func TourismLinepoints(c *gin.Context) {
	lineID := c.PostForm("lineID")
	respData := map[string]interface{}{}
	tlp := []model.TourismLinePoint{}
	err := logic.GetTourismLinePointByLineID(&tlp, lineID)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}

	respData["suc"] = true
	respData["data"] = tlp
	c.JSON(http.StatusOK, respData)
}
