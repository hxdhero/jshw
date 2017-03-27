package logic

import (
	"local/jshw/conf"
	"errors"
	"local/jshw/model"
	"local/jshw/model/sms"
)

//判断用户是否存在,根据传入的结构体中已有的非空数据判断
func IsExist(u *model.User) (bool, error) {
	return conf.DBEngine.Get(u)
}

//根据id查询用户
func FindUserById(user *model.User, id int) error {
	users := []model.User{}
	err := conf.DBEngine.Where(" id = ? ", id).Find(&users)
	if err != nil {
		return err
	}
	if len(users) == 1 {
		user = &users[0]
	} else {
		return errors.New("用户不存在")
	}
	return nil
}

//根据账号获取用户
func FindUserByAccount(user *model.User, account string) error {
	users := []model.User{}
	err := conf.DBEngine.Where(" user_name = ? ", account).Find(&users)
	if err != nil {
		return err
	}
	if len(users) == 1 {
		*user = users[0]
	} else if len(users) > 1 {
		return errors.New("存在重复用户")
	} else {
		return errors.New("用户不存在")
	}
	return nil
}

//根据用户id获取联系人列表
func GetUserContacts(user *[]model.UserContact, userID string) error {
	err := conf.DBEngine.Where(" dt_user_id = ? ", userID).Find(user)
	if err != nil {
		return err
	}
	return nil
}

//添加用户验证码
func AddUserCode(uc *sms.UserCode) error {
	_, err := conf.DBEngine.Insert(uc)
	return err
}

//查询用户验证码
func GetUserCode(uc *sms.UserCode) (b bool, e error) {
	boolean, err := conf.DBEngine.
		Where(" type = ? ", uc.CodeType).
		And(" user_mobile = ? ", uc.UserMobile).
		And("status = ? ", 0).
		Desc("add_time").Limit(1, 0).
		Get(uc)
	return boolean, err
}

//更新用户验证码
func UpdateUserCode(uc *sms.UserCode) error {
	_, err := conf.DBEngine.
		Where(" type = ? ", uc.CodeType).
		And(" str_code = ? ", uc.CodeStr).
		And(" user_mobile = ? ", uc.UserMobile).Update(uc)
	return err
}
