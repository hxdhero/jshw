package main

import (
	"net/url"
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
	"sort"
	"bytes"
)

func main() {
	log.SetFlags(log.Llongfile|log.Ldate)
	alih5new()
}


func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))

}

func alih5new() {
	privateKey := "MIICXAIBAAKBgQDMtmFV5RjuwzmvpdenizGSnR60XF954ynY4t6r9IA7WgeGEFAUfhW5gagbFMpc6FHiUdl29oPjLud+YtWFCkHjobl6bAw0diBAjZeqRBow22RuFzuhrVm3dmVW2vGpAeM25lxSxTSsilwfH+h3I4mV0ySA3t0euIP3ce6rykKIbwIDAQABAoGAUWLEycheJDZ7TaiqVxLQr5BFr8D1uFimv3JawpRfErmVOihsHemOq4Svl6ypU0yNmWOfCFuzTXPNVwLmDpFoZefrcSIw61/6l6R2lCHvk8BNJ40gU573I/UlxNCOL9DiW2ooQIGS+uKV551tFsfG1O8MSZ88KDRoMcqcugPH/4ECQQD0z04J8VUjlYzhSuslS4XajW7dEJjJYWq15osNqaA5PgICkeaR5Pz9Wk8ztElhu0+/dC35cl8hEVQxO5LZikwvAkEA1hHeChGzIQCrQI14YehYQVWgttfqIl0FOvEVzaogjPbfRDhcrZVNxeNhjRFTfN/xTgeWDqf+1PrsZ43Iyi83wQJBANxKUyH1TTSZFU2B2fkUbZ2N2W4JykKka57Flukzc18vIiXn3j/4e4MLqeuP1tyf7hIM3HXz6hBahJVM00b4ALcCQGO1wtSx1dvjceEJhC8miCU2eztvarFC3rLLpLo9Khg+zVP7ZL+9sZIhDUkl7ttVfBI6WlzNR1dw4TiCxCnYwIECQGTZAaBbqPeBiPMe5wfnOACBgposX+lqKRlPIxAE2LaExJw/oZbaBp6TqWcT/7m2vyD3b5gOaAlW+NEjjtr1NeE="
	//base

	//payUrl:="https://openapi.alipaydev.com/gateway.do"//沙箱地址
	//appid := "2016080300158632"//沙箱环境
	appid:="2017040106522852"//正式坏境
	payUrl:="https://openapi.alipay.com/gateway.do"//正式环境
	method := "alipay.trade.wap.pay"
	charset := "utf-8"
	sign_type := "RSA"
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	version := "1.0"
	notifyUrl:="http://www.easytool.site/jshw_notify_ali"
	//biz
	subject := "大乐透"
	out_trade_no := "hxd201704110001"
	total_amount := "0.1"
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
	uv.Set("notify_url",notifyUrl)
	uv.Set("biz_content", biz_content)
	paramKeys:=[]string{}
	for k:=range uv{
		paramKeys=append(paramKeys,k)
	}
	sort.Strings(paramKeys)
	var buf bytes.Buffer
	for _, k := range paramKeys {
		vs := uv[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	bufStr:=buf.String()
	log.Println("待签名: ",bufStr)
	sign, err := Sha1WithRSAPKCS8Base64Sign(bufStr, privateKey)
	if err != nil {
		log.Println(err)
	}
	log.Println("sign:  ", sign)
	uv.Set("sign", sign)
	body := uv.Encode()
	log.Println("body:  ", body)
	uri,err:=url.Parse(payUrl)
	if err != nil {
		log.Println(err)
	}
	uri.RawQuery=body
	log.Println("======  ",uri.String())

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


func Sha1WithRSAPKCS8Base64Sign(data string, privatekey string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(privatekey)
	if err != nil {
		log.Println(err)
		panic("")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		log.Println(err)
		panic("")
	}
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := base64.StdEncoding.EncodeToString(signature)
	return out, nil
}

func Sha1WithRSAPKCS8Base64VerifySign(originData, signData, publickey string) error {
	key, _ := base64.StdEncoding.DecodeString(publickey)
	pub, _ := x509.ParsePKIXPublicKey(key)
	h := sha1.New()
	h.Write([]byte([]byte(originData)))
	hash := h.Sum(nil)
	sig, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash, sig)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return err
	}
	return nil
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
