package model

import "time"

//线路日期表
type TourismLineDate struct {
	ID            int       `xorm:"id" json:"tourismLineDateID"`
	TourismLineID int       `xorm:"dt_dz_tourismline_id"`           //线路表id
	StartDate     time.Time `xorm:"startDate" json:"lineStartDate"` //开始时间
	EndDate       time.Time `xorm:"endDate" json:"lineEndDate"`     //结束时间
	Price1        float64   `xorm:"price1" json:"price"`                         //价格
	Price2        float64   `xorm:"price2" json:"-"`
}

func (t *TourismLineDate) TableName() string {
	return "dt_dz_tourismlinedate"
}
