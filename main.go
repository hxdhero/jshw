package main

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
	"local/jshw/controller"
	"github.com/tommy351/gin-sessions"
	"net/http"
	"local/jshw/conf"
)

/**
返回数据格式:
{
	status: 200/404...  1.先判断通讯是否成功
	suc:    true/false	\2.在判断业务逻辑是否成功
	data:   xxx		3.如果业务逻辑正确,这里返回当前请求的业务数据
	msg:    xxx		4.如果业务逻辑错误,这里返回错误信息,可以显示给用户
}
*/
func main() {
	log.SetFlags(log.Llongfile | log.Ltime)
	time.Local,_=time.LoadLocation("Asia/Shanghai")
	//初始化服务
	conf.InitApp()
	//默认实例
	server:=gin.Default()
	//设置session
	store := sessions.NewCookieStore([]byte("secret123"))
	server.Use(sessions.Middleware("ginsession", store))
	//设置静态路径
	server.Static("/static","./static")
	//设置模板路径
	server.LoadHTMLGlob("view/*")
	server.StaticFile("/favicon.ico","./assets/favicon.ico")
	//设置路由
	server.GET("/",controller.HTMLIndex)//首页
	//----线路相关api
	server.GET("/tourismList",controller.TourismLineList)//获取线路列表
	server.POST("/tourismDetail",controller.TourismLineDetail)//获取线路详情
	server.POST("/tourismPoints",controller.TourismLinepoints)//获取线路集合点
	//----用户相关api
	server.POST("/user_login",controller.LoginJSON)//用户登陆
	server.POST("/user_order_list",controller.OrderListJSON)//获取用户订单列表
	server.POST("/user_contacts",controller.UserContacts)//获取用户联系人
	//----订单相关
	server.POST("/commitOrder",controller.CommitOrder)//提交订单
	server.POST("/orderDetail",controller.OrderDetailByID)//获取订单详情
	server.POST("/orderPay",controller.OrderPay)//订单支付
	//----基础相关api
	server.POST("/sendCodeSMS",controller.SendCodeSMS)//发送短信验证码
	//----支付宝回调
	server.POST("/jshw_notify_ali",controller.AliNotify)//支付宝h5回调
	//*****测试api
	server.POST("/test",controller.TestData)//测试数据
	server.POST("/upload",controller.TestFileUpload)//测试文件上传
	server.GET("/testCss",controller.TestCss)//测试css加载页面
	server.GET("/testCssFile",controller.TestCssFile)//测试css加载页面
	server.GET("/testFlex",controller.TestFlex)

	//匹配前端路由,如果服务器没有匹配到路径就返回index.html给请求方
	server.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	log.Println("server start on:",conf.AppConfig.String("httpport"))
	//启动服务
	server.Run(":" + conf.AppConfig.String("httpport"))

}
