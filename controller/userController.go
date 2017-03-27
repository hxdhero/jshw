package controller

import (
	"github.com/gin-gonic/gin"
	"local/jshw/model"
	"local/jshw/util"
	"net/http"
	"local/jshw/logic"
	"log"
	model_order "local/jshw/model/order"
)

//用户登陆
func LoginJSON(c *gin.Context) {
	respData := map[string]interface{}{}
	account := c.PostForm("account")
	pwd := c.PostForm("pwd")

	user := model.User{}
	err := logic.FindUserByAccount(&user, account)
	log.Println(user)
	if err != nil {
		ReturnError(c, err)
		return
	}
	DESPwd, err := util.DESEncode(pwd, user.Salt)
	if err != nil {
		ReturnError(c, err)
		return
	}

	log.Println(DESPwd)
	if user.Pwd != DESPwd {
		ReturnErrorStr(c, "账号或者密码错误")
		return
	}
	respData["suc"] = true
	respData["data"] = user
	c.JSON(http.StatusOK, respData)
}

//获取用户订单列表
func OrderListJSON(c *gin.Context) {
	respData := map[string]interface{}{}
	userId := c.PostForm("userId")
	if util.IsNull(userId) {
		ReturnErrorStr(c, "数据错误!")
		return
	}

	orderOGs := []model_order.OrderOG{}
	err := logic.OrderGOSummary(&orderOGs, userId)
	if err != nil {
		ReturnError(c, err)
		return
	}
	respData["suc"] = true
	respData["data"] = orderOGs
	c.JSON(http.StatusOK, respData)
}

//获取用户联系人
func UserContacts(c *gin.Context){
	respData := map[string]interface{}{}
	userId := c.PostForm("userId")
	if util.IsNull(userId) {
		log.Println("用户id为空")
		ReturnErrorStr(c, "数据错误!")
		return
	}

	userContacts := []model.UserContact{}
	err:=logic.GetUserContacts(&userContacts,userId)
	if err != nil {
		ReturnError(c, err)
		return
	}
	respData["suc"] = true
	respData["data"] = userContacts
	c.JSON(http.StatusOK, respData)
}