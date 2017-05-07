package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"io/ioutil"
	"encoding/json"
	"local/jshw/model/sms"
	"local/jshw/logic"
	"local/jshw/model/order"
	"strconv"
	"time"
	"fmt"
	model2 "local/jshw/model/tourismline"
	"net/http"
	"local/jshw/conf"
	model3 "local/jshw/model"
	"local/jshw/util"
	"log"
	"net/url"
	"sort"
	"bytes"
)

//用户提交订单
func CommitOrder(c *gin.Context) {
	bodyByte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		seelog.Error(err)
		return
	}
	ma := map[string]interface{}{}
	err = json.Unmarshal([]byte(string(bodyByte)), &ma)
	if err != nil {
		seelog.Error(err)
	}
	seelog.Debug("ma: ", ma)
	seelog.Debug(ma["orderForm"].(map[string]interface{})["orderUserForm"].(map[string]interface{})["mobile"])
	orderData := ma["orderForm"].(map[string]interface{})
	seelog.Debug(orderData)

	lineID := orderData["lineID"].(float64)
	lineDateID:=orderData["lineDateID"].(int)
	//lineIDInt, err := strconv.Atoi(lineID)
	//if err != nil {
	//	seelog.Error(err)
	//	ReturnErrorStr(c, "数据格式错误")
	//	return
	//}
	valCode := orderData["valCode"].(string)
	//订单提交人信息
	committerMobile := orderData["orderUserForm"].(map[string]interface{})["mobile"].(string)
	committerName := orderData["orderUserForm"].(map[string]interface{})["name"].(string)
	committerID := int(orderData["orderUserForm"].(map[string]interface{})["id"].(float64))
	//订单集合点信息
	pointID := int(orderData["selectPoint"].(map[string]interface{})["ID"].(float64))
	pointName := orderData["selectPoint"].(map[string]interface{})["Name"].(string)
	//选择的联系人列表
	selectContacts := orderData["selectContacts"].([]interface{})
	//选择的出行人列表
	selectPersons := orderData["selectPersons"].([]interface{})
	//校验验证码
	uc := sms.UserCode{CodeType: "FillOrder", UserMobile: committerMobile}

	has, err := logic.GetUserCode(&uc)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	if !has {
		seelog.Error("[查询验证码],userCode=", uc)
		ReturnErrorStr(c, "没有找到验证码")
		return
	}
	if uc.CodeStr != valCode {
		seelog.Error("[校验验证码],数据库信息:", uc, "用户提交的验证码: ", valCode)
		ReturnErrorStr(c, "验证码错误")
		return
	}
	//校验验证码是否过期
	expired := conf.AppConfig.String("code_expired")
	expiredInt, err := strconv.Atoi(expired)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	if time.Now().Sub(uc.EffTime).Minutes() > float64(expiredInt) {
		ReturnErrorStr(c, "验证码已经过期")
	}
	//订单信息
	lastOrder := model.Order{}
	_, err = logic.GetLastOrder(&lastOrder)
	if err != nil {
		seelog.Error(err)
		ReturnError(c, err)
		return
	}
	//获取最后一条订单号
	lno := lastOrder.OrderNO
	sno := lno[8:]
	ino, err := strconv.Atoi(sno)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	if ino < 2999 {
		ino++
	} else {
		ino = 1
	}
	//生成订单号
	dStr := time.Now().Format("20060102")
	sino := fmt.Sprintf("%05d", ino)
	seelog.Debug("sino: ", sino)
	orderNOStr := dStr + sino
	seelog.Debug(orderNOStr)
	//获取线路
	line := model2.TourismLine{ID: int(lineID)}
	has, err = logic.GetLineByID(&line)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	if !has {
		seelog.Error("找不到线路: ", lineID)
		ReturnErrorStr(c, "内部错误")
		return
	}
	//获取线路日期
	lineDate:= &model2.TourismLineDate{ID:lineDateID}
	has, err = logic.GetLineDateByID(lineDate)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	if !has {
		seelog.Error("找不到线路日期")
		ReturnErrorStr(c, "找不到线路日期")
		return
	}
	//生成订单
	order := model.Order{
		OrderNO:        orderNOStr,
		UserID:         committerID,
		UserName:       committerName,
		PaymentStatus:  1,
		Mobile:         committerMobile,
		OrderAmount:    line.MaxPrice,
		Status:         1,
		AddTime:        time.Now(),
		PlayNum:        len(selectContacts) + len(selectPersons),
		OutDate:        lineDate.StartDate,
		BackDate:       lineDate.EndDate,
		RendezvousID:   pointID,
		RendezvousName: pointName,
	}
	contacts, err := TranContacts(selectContacts)
	outers, err := TranContacts(selectPersons)
	log.Println("selectContacts: ", contacts)
	log.Println("selectPersons: ", outers)

	//return
<<<<<<< HEAD
	val, err := logic.AddOrder(&order, int(lineID), &contacts, &outers)
=======
	val, err := logic.AddOrder(&order,int(lineID),lineDateID,&contacts,&outers)
>>>>>>> 97c194d2e8f247ef2211abc97549f239c59712f9
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "内部错误")
		return
	}
	if val != 1 {
		seelog.Error("添加order失败")
		ReturnErrorStr(c, "内部错误")
		return
	}

	seelog.Debug("insert val: ", val)
	seelog.Debug("after insert", order)
	c.JSON(http.StatusOK, gin.H{"suc": true})

}

//把用户联系人和出行人封装为对象
func TranContacts(in []interface{}) ([]model3.UserContact, error) {
	out := []model3.UserContact{}
	for _, ele := range in {
		euc := ele.(map[string]interface{})
		sexStr, birthday, err := util.GetIDCardInfo(euc["idcard"].(string))
		if err != nil {
			seelog.Error("身份证信息提取错误", err)
			return nil, err
		}
		uc := model3.UserContact{
			IDCard:        euc["idcard"].(string),
			ContactName:   euc["name"].(string),
			ContactMobile: euc["mobile"].(string),
			Sex:           sexStr,
			Birthday:      birthday,
		}

		out = append(out, uc)
	}
	return out, nil
}

//根据id获取订单详情
func OrderDetailByID(c *gin.Context) {
	orderID := c.PostForm("OrderID")
	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "订单id错误")
		return
	}
	//订单
	order := &model.Order{ID: int64(orderIDInt)}
	err = logic.GetOrderByID(order)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "获取订单详情失败")
		return
	}
	//订单商品
	og, err := logic.GetOrderGoodsByOrderID(order.ID)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, err.Error())
		return
	}

	ols, err := logic.GetOrderLinkersByOrderID(order.ID)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, err.Error())
		return
	}
	data := map[string]interface{}{}
	odm := map[string]interface{}{}
	odm["order"] = order
	odm["ordergoods"] = og
	odm["orderlinkers"] = ols
	data["suc"] = true
	data["data"] = odm
	c.JSON(http.StatusOK, data)
}

//订单支付
func OrderPay(c *gin.Context) {
	data := map[string]interface{}{}
	orderID := c.PostForm("OrderID")
	rUrl := c.PostForm("returnUrl")
	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "订单id错误")
		return
	}
	//订单
	order := &model.Order{ID: int64(orderIDInt)}
	err = logic.GetOrderByID(order)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, "获取订单详情失败")
		return
	}
	order.TradeNO = time.Now().Format("20060102") + strconv.Itoa(int(time.Now().UnixNano()))
	url, err := alih5new(order,rUrl)
	if err != nil {
		seelog.Error(err)
		ReturnErrorStr(c, err.Error())
		return
	}
	data["suc"] = true
	data["data"] = url
	c.JSON(http.StatusOK, data)
}

func alih5new(order *model.Order, returnUrl string) (string, error) {
	privateKey := conf.AppConfig.String("appKey")
	//base
	appid := conf.AppConfig.String("appid")   //正式坏境
	payUrl := conf.AppConfig.String("payUrl") //正式环境
	method := conf.AppConfig.String("payMethod")
	charset := "utf-8"
	sign_type := "RSA"
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	version := "1.0"
	notifyUrl := conf.AppConfig.String("alipay_notify_url")
	//biz
	subject := "订单支付"
	out_trade_no := order.TradeNO
	total_amount := strconv.FormatFloat(order.OrderAmount, 'f', 2, 64)
	product_code := "QUICK_WAP_PAY"

	bizMap := map[string]string{}
	bizMap["subject"] = subject
	bizMap["out_trade_no"] = out_trade_no
	bizMap["total_amount"] = total_amount
	bp:=url.Values{}
	bp.Set("orderID",strconv.Itoa(int(order.ID)))
	bpStr:=bp.Encode()
	bizMap["passback_params"] =url.QueryEscape(bpStr)
	bizMap["product_code"] = product_code
	bizByte, err := json.Marshal(bizMap)
	if err != nil {
		log.Println(err)
		return "", err
	}
	biz_content := string(bizByte)

	uv := url.Values{}
	uv.Set("app_id", appid)
	uv.Set("method", method)
	uv.Set("charset", charset)
	uv.Set("sign_type", sign_type)
	uv.Set("timestamp", timestamp)
	uv.Set("version", version)
	uv.Set("return_url",returnUrl)
	uv.Set("notify_url", notifyUrl)
	uv.Set("biz_content", biz_content)
	paramKeys := []string{}
	for k := range uv {
		paramKeys = append(paramKeys, k)
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
	bufStr := buf.String()
	log.Println("待签名: ", bufStr)
	sign, err := util.Sha1WithRSABase64(bufStr, privateKey)
	if err != nil {
		seelog.Error(err)
		return "", err
	}
	log.Println("sign:  ", sign)
	uv.Set("sign", sign)
	body := uv.Encode()
	log.Println("body:  ", body)
	uri, err := url.Parse(payUrl)
	if err != nil {
		seelog.Error(err)
		return "", err
	}
	uri.RawQuery = body
	seelog.Debug("======  ", uri.String())
	return uri.String(), err

}
