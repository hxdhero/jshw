package model

//线路摘要,列表页显示
type TourismLineList struct {
	tl TourismLine     `json:"line"`
	tld TourismLineDate `json:"lineDate"`
}
