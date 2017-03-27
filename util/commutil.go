package util

import (
	"time"
	"mime/multipart"
	"os"
	"log"
	"bytes"
	"math/rand"
	"sort"
	"crypto/md5"
	"strings"
	"encoding/hex"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/cihub/seelog"
	"runtime/debug"
)

//把时间戳转为yyyy-mMM-dd HH:mm:ss的格式
func FormatUnixTime(ut int64) string {
	sut := time.Unix(ut, 0).Format("2006-01-02 15:04:05")
	return sut
}

/***
获取过去30天的开始和结束日期
*/
func GetStartAndEndTime() (start, end time.Time, err error) {
	time.Local = time.UTC //设置时区
	//startDateStr := "2016-01-01 00:00:00"
	//startDate, err := time.Parse("2006-01-02 15:04:00", startDateStr)
	//return startDateTarget, err

	cd := time.Now()
	startDate := time.Date(cd.Year(), cd.Month(), cd.Day()-30, 0, 0, 0, 0, time.Local)
	endDate := time.Date(cd.Year(), cd.Month(), cd.Day()-1, 23, 59, 59, 0, time.Local)
	return startDate, endDate, nil
}

/**
获取过去30天的日期列表
@param order 1正序  -1倒叙 默认为正序
*/
func GetDaySlice(fmt string) []string {
	days := []string{}
	for i := 30; i > 0; i-- {
		cd := time.Now()
		date := time.Date(cd.Year(), cd.Month(), cd.Day()-i, 0, 0, 0, 0, time.Local)
		if !IsNull(fmt) {
			days = append(days, date.Format(fmt))
		} else {
			days = append(days, date.Format("2006-01-02"))
		}

	}

	return days
}

//获取昨天的日期 yyyy-MM-dd 格式
func GetYesterday(fmt string) string {
	cd := time.Now()
	yd := time.Date(cd.Year(), cd.Month(), cd.Day()-1, 0, 0, 0, 0, time.Local)
	if !IsNull(fmt) {
		return yd.Format(fmt)
	}
	return yd.Format("2006-01-02")
}

// 获取文件大小的接口
type Size interface {
	Size() int64
}
// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}
//获取文件的大小(单位:byte)
func GetFileSize(file multipart.File) int64{
	var fileSize int64
	if statInterface, ok := file.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		log.Println("fileInfo: ",fileInfo)
		log.Println("fileInfoSize: ",fileInfo.Size())
		fileSize = fileInfo.Size()
	}
	if sizeInterface, ok := file.(Size); ok {
		log.Println("siezeInfo: ",sizeInterface.Size())
		fileSize = sizeInterface.Size()
	}
	return fileSize
}

//获取验证码(6位)
func GetValCode() string{
	num:=6
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	strs:="0123456789"
	length:=len(strs)
	var buff bytes.Buffer
	for i:=0;i<num ;i++  {
		random:=r.Intn(length)
		buff.WriteString(strs[random:random+1])
	}
	return buff.String()
}

//获取身份证件信息
func GetIDCardInfo(idcardStr string)(string,time.Time,error){
	//从身份证获取性别
	genderStr:=idcardStr[len(idcardStr)-2:len(idcardStr)-1]
	seelog.Debug("==============================genderStr: ",genderStr)
	if strings.Contains("24680",genderStr){
		genderStr = "女"
	}else{
		genderStr = "男"
	}
	//从身份证获取生日
	birthStr:=idcardStr[6:14]
	birthDay,err:=time.Parse("20060102",birthStr)
	if err != nil {
		seelog.Error(err)
		seelog.Error(string(debug.Stack()))
	}
	return genderStr,birthDay,err
}

//阿里大于发送短信=========================================================================
type AliSMS struct {
	Url       string
	ApiMethod string
	AppKey    string
	Secret    string
	SMSParam  map[string]string
}

//构造阿里大于短信发送结构体
func NewAliSMS(url, apiMethod, appKey, secret string) *AliSMS {
	alisms := &AliSMS{
		Url:    url,
		Secret: secret,
	}
	alisms.SMSParam = make(map[string]string)
	alisms.SMSParam["app_key"] = appKey
	alisms.SMSParam["method"] = apiMethod
	alisms.SMSParam["sign_method"] = "md5"
	alisms.SMSParam["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	alisms.SMSParam["format"] = "json"
	alisms.SMSParam["v"] = "2.0"
	//提交的时候添加sing参数

	return alisms
}

//发送短信请求,返回请求响应消息,error
//响应消息失败:{"error_response":{"code":28,"msg":"Missing app key","request_id":"zt9tg1f4lhkg"}}
//成功:{"alibaba_aliqin_fc_sms_num_send_response":{"result":{"err_code":"0","model":"106299406230^1108573321309","success":true},"request_id":"z24ifvpric21"}}
func (a *AliSMS) SendSMS() (suc bool, respStr string, err error) {
	//排序参数中的key
	paramKey := []string{}
	for k, _ := range a.SMSParam {
		paramKey = append(paramKey, k)
	}
	sort.Strings(paramKey)
	//生成signStr,把参数中的key和value拼接起来并且首尾加上秘钥: secret+key1value1key2value2+secret
	var buffer bytes.Buffer
	buffer.WriteString(a.Secret)
	for _, str := range paramKey {
		buffer.WriteString(str)
		buffer.WriteString(a.SMSParam[str])
	}
	buffer.WriteString(a.Secret)
	//生成的signStr
	signStr := buffer.String()
	//用md5加密signStr,生成最终的sign
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStr))
	cipherStr := md5Ctx.Sum(nil)
	signStr = strings.ToUpper(hex.EncodeToString(cipherStr))
	//最终生成的signStr
	a.SMSParam["sign"] = signStr
	//生成请求地址+参数
	smsParam := url.Values{}
	for k, v := range a.SMSParam {
		smsParam.Set(k, v)
	}
	var reqUrlBuffer bytes.Buffer
	reqUrlBuffer.WriteString(a.Url)
	reqUrlBuffer.WriteString("?")
	reqUrlBuffer.WriteString(smsParam.Encode())
	reqUrl := reqUrlBuffer.String()
	seelog.Debug(reqUrl)

	resp, err := http.Get(reqUrl)
	if err != nil {
		seelog.Error(err)
		return false, "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		seelog.Error(err)
		return false, "", err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		seelog.Error(err)
		return false, "", err
	}

	log.Println(string(respBody))
	if resp.StatusCode != 200 {
		return false, string(respBody), err
	}
	if _, ok := data["error_response"]; ok {
		return false, string(respBody), err
	}
	return true, string(respBody), err

}

//=========================================================================