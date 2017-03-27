package model

//线路安排
type TourismLineTrip struct {
	ID            int    `xorm:"id" json:"tourismLineDateID"`
	TourismLineID int    `xorm:"dt_dz_tourismline_id"` //线路id
	Title         string `xorm:"Title"`                //安排标题
	Contents      string `xorm:"contents"`             //安排内容

}

func (t *TourismLineTrip) TableName() string {
	return "dt_dz_tourismlinetrip"
}
