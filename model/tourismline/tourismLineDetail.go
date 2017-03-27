package model

//1对多的线路
type TourismLineDetail struct {
	Line TourismLine `json:"line"`
	LineDate TourismLineDate `json:"lineDate"` //按照date倒叙拿最近的一条
	LineTrip []TourismLineTrip `json:"lineTrip"` //根据天数,一条线路可能对应多条线路
	LineImage []TourismLineImages `json:"lineImages"` //一条线路对应多个图片
	LinePoint []TourismLinePoint `json:"linePoints"` //一条线路对应多个集合点
}