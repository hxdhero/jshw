package model

//1对1的重复线路
type TourismLineDetailAll struct {
	TourismLine       `xorm:"extends"` //线路
	TourismLineDate   `xorm:"extends"` //线路日期
	tlt []TourismLineTrip   `xorm:"extends" json:"tlt"` //线路安排
	til []TourismLineImages `xorm:"extends" json:"til"` //线路图片
}
