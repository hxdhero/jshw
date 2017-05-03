package logic

import (
	"local/jshw/model/order"
	"local/jshw/conf"
	"github.com/cihub/seelog"
	model2 "local/jshw/model/tourismline"
	"errors"
	"strings"
	"log"
	"time"
	"github.com/go-xorm/xorm"
	model3 "local/jshw/model"
)

//根据用户id查询订单和订单商品
func OrderGOSummary(ogos *[]model.OrderOG, userID string) error {
	err := conf.DBEngine.Table("dt_orders").Alias("od").
		Join("INNER", []string{"dt_order_goods", "og"}, "od.id = og.order_id").
		Where(" od.user_id = ? ", userID).Desc("od.add_time").Find(ogos)
	return err
}

//插入订单
//1.保存订单
//2.保存订单商品
//3.保存订单联系人
//4.保存用户联系人
func AddOrder(order *model.Order,lineID,lineDateID int,contacts,outers *[]model3.UserContact) (int64,error) {
	//开启事务
	session:=conf.DBEngine.NewSession()
	defer session.Close()
	err:=session.Begin()
	if err != nil {
		seelog.Error(err)
		return 0,err
	}
	//1.保存订单
	val,err:=session.Insert(order)
	if val != 1{
		session.Rollback()
		return val,err
	}
	if err != nil {
		seelog.Error(err)
		session.Rollback()
		return 0,err
	}
	//线路
	line:=model2.TourismLine{ID:lineID}
	has,err :=GetLineByID(&line)
	if !has{
		seelog.Error("找不到订单: ",lineID)
		return 0,errors.New("内部错误")
	}
	if err != nil {
		seelog.Error(err)
		return 0,errors.New("内部错误")
	}
	//线路日期
	//lineDate:=model2.TourismLineDate{TourismLineID:int(lineID)}
	//has,err=GetTourismLineDateByLineID(&lineDate)
	//if !has{
	//	seelog.Error("找不到线路日期: ",lineID)
	//	return 0,errors.New("内部错误")
	//}
	//if err != nil {
	//	seelog.Error(err)
	//	return 0,errors.New("内部错误")
	//}
	//2.保存订单商品
	orderGoods:=model.OrderGoods{
		OrderID:order.ID,
		GoodsTitle:line.Title,
		GoodsPrice:line.MaxPrice,
		RealPrice:line.MaxPrice,
		Quantity:1,
		TourismLineID:int64(lineID),
		TourismLineDateID:int64(lineDateID),
	}
	val,err=session.Insert(&orderGoods)
	if val!=1{
		seelog.Error("保存ordergoods失败,orderID:",order.ID,"  lineTitle: ",line.Title,"error: ",err)
		session.Rollback()
		return 0,errors.New("内部错误")
	}
	if err != nil {
		seelog.Error(err)
		session.Rollback()
		return 0,err
	}
	//3.保存订单成员
	_,err=session.Insert(contacts)
	if err != nil {
		session.Rollback()
		return 0,err
	}
	session.Commit()
	session.Close()
	//返回订单保存结果
	return val,err
}

//Deprecated
func saveOrderLinkers(session xorm.Session,orderID int64,eleMap map[string]interface{}) (int64,error){
	idcardStr:=eleMap["idcard"].(string)
	//从身份证获取性别
	genderStr:=idcardStr[len(idcardStr)-2:len(idcardStr)-1]
	if strings.Contains("24680",genderStr){
		genderStr = "男"
	}else{
		genderStr = "女"
	}
	//从身份证获取生日
	birthStr:=idcardStr[6:14]
	birthDay,err:=time.Parse("20060102",birthStr)
	if err != nil {
		log.Println(err)
		return 0,err
	}
	orderLinker:=model.OrderLinker{
		OrderID:orderID,
		Name:eleMap["name"].(string),
		Mobile:eleMap["mobile"].(string),
		IDCard:idcardStr,
		Gender:genderStr,
		BirthDay:birthDay,
	}
	return session.Insert(&orderLinker)
}

//获取最后提交的一条订单(为了得到订单号)
func GetLastOrder(order *model.Order) (bool,error){
	return conf.DBEngine.Desc("order_no").Get(order)
}
