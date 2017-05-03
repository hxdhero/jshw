package model

//参加线路的用户
type LineUser struct {
	Mobile int    `xorm:"mobile"`
	Name   string `xorm:"name"`
	Sex    string `xorm:"sex"`
}

func (u *LineUser) TableName() string {
	return "dt_order_linkers"
}
