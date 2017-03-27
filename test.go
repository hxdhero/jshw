package main
/*

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

func main() {

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	m2()
}

func m2(){
	user:=&User{}
	setUser(user)
	log.Println("3.==",user)
}

func setUser(u *User){
	uu:=User{Id:11,UserName:"hero"}
	*u=uu
	log.Println("1.==",uu)
	log.Println("2.==",u)
}

func m1(){
	dburl := "odbc:server=139.196.187.30;user id=jshwclub.com;password={easy_9999};database=jshwclub"

	log.Println(dburl)
	DBEngine, err := xorm.NewEngine("mssql", dburl)
	if err != nil {
		log.Println(err)
	}
	DBEngine.ShowSQL(true)
	DBEngine.Logger().SetLevel(core.LOG_DEBUG)
	err = DBEngine.Sync2(new(User))
	if err != nil {
		log.Println(err)
	}
	users := []User{}
	err = DBEngine.Where(" user_name = ? ", "13761764141").Find(&users)
	if err != nil {
		log.Println(err)
		return
	}
	count, err := DBEngine.Count(&User{})
	log.Println("count :", count)
	log.Println(users)
	log.Println(users[0])
	//
	user := User{UserName: "13761764141", Pwd: "E0602A202F6EFB9B"}
	isExit, err := DBEngine.Get(&user)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(isExit)
	//
	uu := &[]User{}
	err = DBEngine.Where(" user_name = ? ", "13761764141").Find(uu)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(uu)
}

type User struct {
	Id           int64     `xorm:"id"`            //id
	GroupID      int64     `xorm:"group_id"`      //用户组 默认0
	UserName     string    `xorm:"user_name"`     //账号
	Pwd          string    `xorm:"password"`      //密码
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
*/
