package model

//线路集合地点
type TourismLinePoint struct {
	ID     int    `xorm:"id"` //id
	Name   string `xorm:"name"` //集合点名称
	Orders int    `xorm:"orders"` //排序
}

func (t *TourismLinePoint)TableName() string {
	return "dt_dz_point"
}