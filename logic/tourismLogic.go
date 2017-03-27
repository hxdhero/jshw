package logic

import (
	"local/jshw/conf"
	"local/jshw/model/tourismline"
	"github.com/cihub/seelog"
	"strings"
)

//获取开放的路线,根据orders排序
func GetTourismList(tls *[]model.TourismLine) error {
	err := conf.DBEngine.
		Where(" is_top = ? ", 1).
		Asc(" orders ").
		Cols("id","code","title","images","minPrice","maxPrice").
		Find(tls)
	return err
}

//查询线路和相关数据
func GetTourismDetail(tld *[]model.TourismLineDetailAll, tourismLineID int) error {
	err := conf.DBEngine.Table("dt_dz_tourismline").Alias("line").
		Join("INNER", []string{"dt_dz_tourismlinedate", "lineDate"}, "line.id = lineDate.dt_dz_tourismline_id").
		Join("INNER", []string{"dt_dz_tourismlinetrip", "lineTrip"}, "line.id = lineTrip.dt_dz_tourismline_id").
		Join("INNER", []string{"dt_dz_tourismlineimages", "lineImage"}, "line.id = lineImage.dt_dz_tourismline_id").
		Where("line.id = ? ", tourismLineID).
		Desc("lineDate.startDate").
		Find(tld)
	return err
}

//根据id获取线路
func GetLineByID(tl *model.TourismLine) (bool,error){
	return conf.DBEngine.Where(" id = ? ",tl.ID).Get(tl)
}

//根据id查询线路
func GetTourismLineByID(tl *[]model.TourismLine, tourismLineID int) error {
	err := conf.DBEngine.Where(" id = ? ", tourismLineID).Find(tl)
	return err
}

//根据线路查找日期
func GetTourismLineDateByLineID(tld *model.TourismLineDate) (bool,error) {
	return conf.DBEngine.Where(" dt_dz_tourismline_id = ? ", tld.TourismLineID).Desc("startDate").Limit(1).Get(tld)
}

//根据线路查找路线
func GetTourismLineTripByLineID(tlt *[]model.TourismLineTrip, tourismLineID int) error {
	err := conf.DBEngine.Where(" dt_dz_tourismline_id = ? ", tourismLineID).Find(tlt)
	return err
}

//根据线路查找图片
func GetTourismLineImagesByLineID(tli *[]model.TourismLineImages, tourismLineID int) error {
	err := conf.DBEngine.Where(" dt_dz_tourismline_id = ? ", tourismLineID).Find(tli)
	return err
}

//根据线路查询集合点
func GetTourismLinePointByID(tlp *[]model.TourismLinePoint, pointID []string) error {
	err := conf.DBEngine.In(" id ", pointID).Find(tlp)
	return err
}

//根据线路id获取集合点
func GetTourismLinePointByLineID(tlp *[]model.TourismLinePoint, lineID string) error {
	lines := []model.TourismLine{}
	err := conf.DBEngine.Where(" id = ? ", lineID).Find(&lines)
	if err != nil {
		seelog.Error(err)
	}
	err = conf.DBEngine.In(" id ", strings.Split(lines[0].Rendezvous, ",")).Find(tlp)
	if err != nil {
		seelog.Error(err)
	}
	return err
}
