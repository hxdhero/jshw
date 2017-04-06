package main

import (
	"net/url"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"
	"encoding/json"
	"crypto/x509"
	"crypto/sha1"
	"crypto/rsa"
	"crypto"
	"fmt"
	"crypto/rand"
	"encoding/base64"
)

func wx() {
	appID := "wx4799ba69be0a41ad"
	mchID := "1275229801"
	nonce_str := "5K8264ILTKCH16CQ2502SI8ZNMTM67VS"
	body := "线路详情-线路预订"
	out_trade_no := "20150806125346"
	total_fee := "0.5"
	spbill_create_ip := "123.12.12.123"
	notify_url := "http://www.easytool.site"
	trade_type := "MWEB"
	scene_info := `{"h5_info": {"type":"Wap","wap_url": "https://www.easytool.site","wap_name": "金山户外"}}`
	key := "55014f2906d9f628124874497e8fb16f"
	//sign:=""

	uv := url.Values{}
	uv.Set("appid", appID)
	uv.Set("mch_id", mchID)
	uv.Set("nonce_str", nonce_str)
	uv.Set("body", body)
	uv.Set("out_trade_no", out_trade_no)
	uv.Set("total_fee", total_fee)
	uv.Set("spbill_create_ip", spbill_create_ip)
	uv.Set("notify_url", notify_url)
	uv.Set("trade_type", trade_type)
	uv.Set("scene_info", scene_info)
	//log.Println(uv)
	//keys := []string{}
	//for k, _ := range uv {
	//	keys = append(keys, k)
	//}
	//sort.Strings(keys)
	//for _, elem := range keys {
	//	log.Println(elem)
	//	uv.Set(elem, uv.Get(elem))
	//}
	paramStr := uv.Encode()
	log.Println(paramStr)
	signStr := paramStr + "&key=" + key
	sign := GetMd5String(signStr)
	log.Println(strings.ToUpper(sign))

}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

func alih5new() {
	privateKey := "MIICXgIBAAKBgQCi+NGFyyhFnMl39Crxc0Yomoen1rkjNf7SQhp0ec5554BXk98xLc8HeklR6kSlpnU0RQ+v2awolEMMFryqZHbz8V9wF6q8vAF2dKgAYCyqDyJZjxOKCiE7H421kC2vo4VQHEXzi6OHuAFD4RLnuXIJuswXAZ/2Jt4kDx5qrAMxsQIDAQABAoGBAKCI83OVDMWNzVOxHIAdajXjCtAFDvglXy9k2ER2HDMvHNioHAqIslAOYJ0lZJu8XeWwReSWSiTq7yTAXPaH4jeVlhk+TwCtRNq4l8akXzAv6H6ztWG2mUePkxxU61CNnVVH6QxpESaKdTTFrum9w1A3S+FR65GSCuUyTKv9vellAkEA09dyJILVaV8ALrJG5mwLZl4uYP8j6ALgX7KmOFnit4mx7rszcVgJznL9k9kiNXFV3+3QnME+ufTTq+ilqwce2wJBAMTxh7idyijWFbMx959EQH2MNZGa7Ciuda39D3ZAgShQc7Tg+omlxIoKQTLSG6jCNeSNNcR8fEcbRjUVgVPCuWMCQQC+oW7eujmPo+THILi6m9m6WeBEevR10TjWBS6dIQ3q+eb7nMwTIBVbCZF1XXzyOLX9V8VVenSW5GEinq2OdU7nAkBG1F7thMI6IZS4V9Yoz5EqFg0GCuO4VdY49vRioRxSdWzHtsokSxv+UWXVcz9DWGWthyO5QNQpdqOvX8ada0DlAkEAqag+nbgIY1XdPzvRegCfOKmnecf+B2gW2Auq5JqVIqBww5RSw2vJOfQ5k2eYFHJcvpchPTQS0hyf43kvijY4jg=="
	//base
	appid := "2017040106522852"
	method := "alipay.trade.wap.pay"
	charset := "utf-8"
	sign_type := "RSA"
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	version := "1.0"

	//biz
	subject := "大乐透"
	out_trade_no := "70501111111S001111119"
	total_amount := "0.5"
	product_code := "QUICK_WAP_PAY"

	bizMap := map[string]string{}
	bizMap["subject"] = subject
	bizMap["out_trade_no"] = out_trade_no
	bizMap["total_amount"] = total_amount
	bizMap["product_code"] = product_code
	bizByte, err := json.Marshal(bizMap)
	if err != nil {
		log.Println(err)
	}
	biz_content := string(bizByte)

	uv := url.Values{}
	uv.Set("app_id", appid)
	uv.Set("method", method)
	uv.Set("charset", charset)
	uv.Set("sign_type", sign_type)
	uv.Set("timestamp", timestamp)
	uv.Set("version", version)
	uv.Set("biz_content", biz_content)
	uvStr := uv.Encode()
	sign, err := Sha1WithRSABase64(uvStr, privateKey)
	if err != nil {
		log.Println(err)
	}
	log.Println("sign:  ", sign)
	uv.Set("sign", sign)
	body := uv.Encode()
	log.Println("body:  ", body)
}


func Sha1WithRSABase64(data string, privatekey string) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(privatekey)
	privateKey, _ := x509.ParsePKCS1PrivateKey(key)
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := base64.StdEncoding.EncodeToString(signature)
	return out, nil
}


func alih5old() {
	md5key := "p8h3q7ymgzmzh6zxsiqplr772tk5ssd1"

	//basicParam
	service := "alipay.wap.create.direct.pay.by.user"
	partner := "2088811282157771"
	_input_charset := "UTF-8"
	sign_type := "MD5"

	//bizParam
	out_trade_no := "70501111111S001111119"
	subject := "户外线路"
	total_fee := "0.5"
	seller_id := "2088811282157771"
	showUrl := "http://www.taobao.com/product/113714.html"

	uv := url.Values{}
	uv.Set("service", service)
	uv.Set("partner", partner)
	uv.Set("_input_charset", _input_charset)
	uv.Set("out_trade_no", out_trade_no)
	uv.Set("subject", subject)
	uv.Set("total_fee", total_fee)
	uv.Set("seller_id", seller_id)
	uv.Set("payment_type", "1")
	uv.Set("show_url", showUrl)
	uvStr := uv.Encode()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(uvStr))
	cipherStr := md5Ctx.Sum([]byte(md5key))
	sign := hex.EncodeToString(cipherStr)
	log.Println("sign: ", string(sign))
	uv.Set("sign_type", sign_type)
	uv.Set("sign", string(sign))
	log.Println(uv.Encode())
}

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
