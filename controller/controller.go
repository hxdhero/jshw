package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"log"
	"os"
	"io"
	"local/jshw/util"
	"local/jshw/conf"
	"encoding/json"
	"github.com/cihub/seelog"
	"local/jshw/logic"
	"local/jshw/model/sms"
	"strconv"
	"sort"
	"strings"
	"encoding/base64"
	"net/url"
	"local/jshw/model/order"
)

//返回错误信息
func ReturnError(c *gin.Context, errorMsg error) {
	log.Println(errorMsg)
	respData := map[string]interface{}{}
	respData["suc"] = false
	respData["msg"] = errorMsg.Error()
	c.JSON(http.StatusOK, respData)
}

//返回错误信息
func ReturnErrorStr(c *gin.Context, errorMsg string) {
	log.Println(errorMsg)
	respData := map[string]interface{}{}
	respData["suc"] = false
	respData["msg"] = errorMsg
	c.JSON(http.StatusOK, respData)
}

//首页
func HTMLIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

//发送短信验证码
func SendCodeSMS(c *gin.Context) {
	seelog.Flush() //打印日志
	respData := map[string]interface{}{}
	mobile := c.PostForm("mobile") //接收短信的手机号
	smsType := c.PostForm("type")  //短信类型:1提交订单验证码
	seelog.Debug("mobile: ", mobile)
	seelog.Debug("smstype: ", smsType)
	//校验手机
	if util.IsNull(mobile) || len(mobile) < 11 {
		ReturnErrorStr(c, "手机号错误")
		return
	}
	//查询该手机号该类型的验证码是否存在
	uc := sms.UserCode{}
	uc.UserMobile = mobile
	switch smsType {
	case "1":
		uc.CodeType = "FillOrder"
	default:
		seelog.Error("发送短信请求,smsType: ", smsType)
		ReturnErrorStr(c, "错误的参数")
		return
	}
	//查询用户验证码 手机号,验证码类型
	_, err := logic.GetUserCode(&uc)
	if err != nil {
		seelog.Error(err)
		return
	}
	expired := conf.AppConfig.String("code_expired")
	expiredInt, err := strconv.Atoi(expired)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	//如果验证码存在则判断时间是否过期
	seelog.Debug("uc", uc)
	if !util.IsNull(uc.CodeStr) {
		if (time.Now().Sub(uc.EffTime).Minutes() <= float64(expiredInt)) {
			seelog.Debug("time: ", time.Now().Sub(uc.EffTime).Minutes())
			seelog.Debug("expired: ", expiredInt)
			ReturnErrorStr(c, "您的验证码还未过期")
			return
		}

	}
	//如果验证码不存在或者已经过期,发送短信并且保存/更新验证码
	valCode := util.GetValCode()
	alisms := util.NewAliSMS(
		conf.AppConfig.String("aldy_url"),
		conf.AppConfig.String("aldy_api_sms_send"),
		conf.AppConfig.String("aldy_appkey"),
		conf.AppConfig.String("aldy_secret"),
	)
	alisms.SMSParam["sms_type"] = "normal"
	alisms.SMSParam["sms_free_sign_name"] = conf.AppConfig.String("aldy_signature")
	alisms.SMSParam["rec_num"] = mobile
	alisms.SMSParam["sms_template_code"] = conf.AppConfig.String("commit_order_code")
	smsParam := map[string]interface{}{}
	switch smsType {
	case "1":
		smsParam["code"] = valCode                                //验证码
		smsParam["valid"] = conf.AppConfig.String("code_expired") //有效期
	default:
		seelog.Error("发送短信请求,smsType: ", smsType)
		ReturnErrorStr(c, "错误的参数")
		return
	}
	jsonStr, err := json.Marshal(smsParam)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "短信发送失败,请拨打客服电话")
		return
	}
	alisms.SMSParam["sms_param"] = string(jsonStr)
	suc, resp, err := alisms.SendSMS()
	seelog.Debug("suc: ", suc, "resp: ", resp, "err: ", err)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "短信发送失败,请拨打客服电话")
		return
	}
	if suc {
		uc.AddTime = time.Now()
		uc.EffTime = time.Now().Add(time.Duration(expiredInt) * time.Minute)
		uc.CodeStr = valCode
		err = logic.AddUserCode(&uc)
		if err != nil {
			seelog.Error(err)
			ReturnErrorStr(c, "内部错误")
			return
		}
		respData["suc"] = true
		respData["data"] = "suc"
		c.JSON(http.StatusOK, respData)
	} else {
		json.Unmarshal([]byte(resp), &respData)
		seelog.Error(respData["error_response"].(map[string]interface{})["msg"])
		ReturnErrorStr(c, respData["error_response"].(map[string]interface{})["msg"].(string))
		return
	}
}

//支付宝回调接口
func AliNotify(c *gin.Context) {
	defer func() {
		seelog.Flush()
		recover()
	}()
	seelog.Flush() //打印日志
	//dump, _ := httputil.DumpRequest(c.Request, true)
	//seelog.Debug("request: ", string(dump))
	c.Request.ParseForm()
	urlValues := c.Request.Form
	//支付宝给回的签名
	sign := urlValues["sign"][0]
	signArray := []string{}
	for k := range urlValues {
		if k != "sign" && k != "sign_type" {
			signArray = append(signArray, k)
		}
	}
	sort.Strings(signArray)

	signParams := map[string]string{}
	for _, elem := range signArray {
		signParams[elem] = urlValues[elem][0]
	}

	strArray := []string{}
	for k, v := range signParams {
		strArray = append(strArray, k+"="+v)
	}
	sort.Strings(strArray)
	//得到待验签字符串
	signStr := strings.Join(strArray, "&")
	seelog.Debug("signStr: ", signStr)
	seelog.Debug("base64SignStr: ", base64.StdEncoding.EncodeToString([]byte(signStr)))
	seelog.Debug("sign: ", sign)
	err := util.Sha1WithRSAPKCS8Base64VerifySign(signStr, sign, conf.AppConfig.String("alipubKey"))
	if err != nil {
		seelog.Error(err)
		panic(err)
	}

	tradeNo := urlValues["out_trade_no"][0]      //订单号
	amount := urlValues["total_amount"][0]       //金额
	sellerID := urlValues["seller_id"][0]        //卖家
	appID := urlValues["app_id"][0]              //应用id
	paramBack := urlValues["passback_params"][0] //返回的透传参数
	seelog.Debug("paramBack", paramBack)
	seelog.Debug("tradeNo: ", tradeNo)
	seelog.Debug("amount: ", amount)
	seelog.Debug("sellerID: ", sellerID)
	seelog.Debug("appID: ", appID)
	//todo 校验以上4个参数
	tradeStatus := urlValues["trade_status"][0] //交易状态
	seelog.Debug(tradeStatus)
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		//订单完成

		ucStr, err := url.QueryUnescape(paramBack)
		if err != nil {
			seelog.Error(err)
		}
		urlBackParam, err := url.ParseQuery(ucStr)
		if err != nil {
			seelog.Error(err)
		}
		orderID := urlBackParam.Get("orderID")
		orderIDInt, err := strconv.Atoi(orderID)
		if err != nil {
			seelog.Error(orderIDInt)
		}
		order := &model.Order{ID: int64(orderIDInt)}
		err = logic.GetOrderByID(order)
		seelog.Debug("获取订单: ", order)
		if err != nil {
			seelog.Error("获取订单失败: ", err)
		}
		amount64, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			seelog.Error(err)
		}
		order.PaymentFee = amount64
		order.PaymentTime = time.Now()
		order.PaymentStatus = 2
		order.Status = 2
		err = logic.UpdateOrder(order)
		if err != nil {
			seelog.Error(err)
		}

	}
	c.Writer.WriteString("success")
	seelog.Debug("ali_notify_suc")
}

//测试方法---------------------
func TestData(c *gin.Context) {
	time.Sleep(time.Second * 1)
	c.JSON(http.StatusOK, gin.H{"name": "bruce", "gender": "male"})
}

func TestCss(c *gin.Context) {
	c.HTML(http.StatusOK, "tt.html", gin.H{})
}

func TestCssFile(c *gin.Context) {
	time.Sleep(3 * time.Second)
	//c.File("./static/css/app.82ba78cc15eca09422469082ea3e00ed.css")
}

func TestFlex(c *gin.Context) {
	c.HTML(http.StatusOK, "flex.html", gin.H{})
}

func TestFileUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("certFile")
	if err != nil {
		log.Println(err)
		return
	}
	fileName := header.Filename

	timeStr := time.Now().Format("2006-01-02")
	savePath := "/home/bruce/桌面/fileUpload/" + timeStr
	err = os.MkdirAll(savePath, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	out, err := os.Create(savePath + "/" + fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(timeStr)
	log.Println("fileName: ", fileName)

	c.JSON(http.StatusOK, gin.H{"fileSize": util.GetFileSize(file)})
	return
}

//测试方法---------------------
