package model

import (
	"time"
)

type User struct {
	Id           int64     `xorm:"id"`            //id
	GroupID      int64     `xorm:"group_id"`      //用户组 默认0
	UserName     string    `xorm:"user_name"`     //账号
	Pwd          string    `xorm:"password" json:"-"`      //密码
	Salt         string    `xorm:"salt"`          //6位随机字符串,加密密码用
	NickName     string    `xorm:"nick_name"`     //昵称
	Avatar       string    `xorm:"avatar"`        //头像
	Email        string    `xorm:"email"`         //邮箱
	Gender       string    `xorm:"sex"`           //性别
	Birthday     time.Time `xorm:"birthday"`      //生日
	TelPhone     string    `xorm:"telphone"`      //联系电话
	Mobile       string    `xorm:"mobile"`        //手机
	QQ           string    `xorm:"qq"`            //qq
	Address      string    `xorm:"address"`       //地址
	SafeQuestion string    `xorm:"safe_question"` //安全问题
	SafeAnswer   string    `xorm:"safe_answer"`   //安全问题答案
	Amount       int64     `xorm:"amount"`        //预存款
	Point        int64     `xorm:"point"`         //积分
	Exp          int64     `xorm:"exp"`           //经验值
	Status       int64     `xorm:"status"`        //用户状态 0.正常 1.待验证 2.待审核
	RegTime      time.Time `xorm:"reg_time"`      //注册时间
	RegIP        string    `xorm:"reg_ip"`        //注册ip
	Level        int64     `xorm:"level"`         //等级
	RealName     string    `xorm:"realname"`      //真实姓名
	PlayCount    int64     `xorm:"playcount"`     //出行次数
	Isbuiltin    int       `xorm:"isbuiltin"`
}

func (u *User) TableName() string {
	return "dt_users"
}
