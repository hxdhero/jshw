package model

import "time"

type OrderLinker struct {
	ID int `xorm:" 'id' pk autoincr <- "`
	OrderID int64 `xorm:"dt_orders_id"`
	Name string `xorm:"name"`
	Mobile string `xorm:"mobile"`
	IDCard string `xorm:"idcard"`
	Gender string `xorm:"sex"`
	BirthDay time.Time `xorm:"birthday"`
	Remark string `xorm:"remark"`
}

func (ol *OrderLinker) TableName() string{
	return "dt_order_linkers"
}