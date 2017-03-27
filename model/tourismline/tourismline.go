package model

import "time"

//线路表
type TourismLine struct {
	ID           int       `xorm:"id" json:"tourismLineID"`
	Code         string    `xorm:"code"`
	Title        string    `xorm:"title"`
	Rendezvous   string    `xorm:"rendezvous"`   //集合地
	Days         int       `xorm:"days"`         //天数
	PlaceExplain string    `xorm:"placeExplain"` //集合地点详细说明
	Remark       string    `xorm:"remarks"`
	CostExplain  string    `xorm:"costExplain"` //费用说明
	Images       string    `xorm:"images"`      //线路图片
	Types        string    `xorm:"types"`       //线路类型编码
	TypeNames    string    `xorm:"typenames"`   //线路类型 自由行,徒步等
	MinPrice     float64   `xorm:"minPrice"`    //最低价格
	MaxPrice     float64   `xorm:"maxPrice"`    //最高价格
	CreateDate   time.Time `xorm:"createDate"`
	UpdateDate   time.Time `xorm:"updateDate"`
	PlaceName    string    `xorm:"placeNames"` //线路位置
}

func (t *TourismLine) TableName() string {
	return "dt_dz_tourismline"
}
