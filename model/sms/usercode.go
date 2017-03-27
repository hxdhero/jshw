package sms

import "time"
//用户验证码表
type UserCode struct {
	ID int `xorm:" 'id' pk autoincr <-"`
	UserID int `xorm:"-"`//用户id
	UserName string `xorm:"-"`//用户名
	CodeType string `xorm:"type"`//验证类型 FillOrder.提交订单
	CodeStr string `xorm:"str_code"`//验证码
	EffTime time.Time `xorm:"eff_time"`//到期时间
	AddTime time.Time `xorm:"add_time"`//添加时间
	UserMobile string `xorm:"user_mobile"`//用户手机号
}

func (u *UserCode) TableName() string{
	return "dt_user_code"
}
