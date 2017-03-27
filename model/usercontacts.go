package model

import "time"

//用户联系人
type UserContact struct {
	ID int `xorm:" 'id' pk autoincr <- "`
	UserID int `xorm:"dt_user_id"`
	ContactName string `xorm:"name"`
	ContactMobile string `xorm:"mobile"`
	IDCard string `xorm:"idcard"`
	Sex string `xorm:"sex"`
	Birthday time.Time `xorm:"birthday"`
	CreateDate time.Time `xorm:"createDate"`
	UpdateDate time.Time `xorm:"updateDate"`
}

func (uc *UserContact)TableName() string{
	return "dt_user_linkers"
}
