package model

//中间表
type OrderOG struct {
	Order `xorm:"extends"`
	OrderGoods `xorm:"extends"`
}
