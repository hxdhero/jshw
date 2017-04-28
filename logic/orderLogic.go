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
func AddOrder(order *model.Order, lineID int, contacts, outers *[]model3.UserContact) (int64, error) {
	//开启事务
	session := conf.DBEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		seelog.Error(err)
		return 0, err
	}
	//1.保存订单
	val, err := session.Insert(order)
	if val != 1 {
		session.Rollback()
		return val, err
	}
	if err != nil {
		seelog.Error(err)
		session.Rollback()
		return 0, err
	}
	//线路
	line := model2.TourismLine{ID: lineID}
	has, err := GetLineByID(&line)
	if !has {
		seelog.Error("找不到订单: ", lineID)
		return 0, errors.New("内部错误")
	}
	if err != nil {
		seelog.Error(err)
		return 0, errors.New("内部错误")
	}
	//线路日期
	lineDate := model2.TourismLineDate{TourismLineID: int(lineID)}
	has, err = GetTourismLineDateByLineID(&lineDate)
	if !has {
		seelog.Error("找不到线路日期: ", lineID)
		return 0, errors.New("内部错误")
	}
	if err != nil {
		seelog.Error(err)
		return 0, errors.New("内部错误")
	}
	//2.保存订单商品
	orderGoods := model.OrderGoods{
		OrderID:           order.ID,
		GoodsTitle:        line.Title,
		GoodsPrice:        line.MaxPrice,
		RealPrice:         line.MaxPrice,
		Quantity:          1,
		TourismLineID:     int64(lineID),
		TourismLineDateID: int64(lineDate.ID),
	}
	val, err = session.Insert(&orderGoods)
	if val != 1 {
		seelog.Error("保存ordergoods失败,orderID:", order.ID, "  lineTitle: ", line.Title, "error: ", err)
		session.Rollback()
		return 0, errors.New("内部错误")
	}
	if err != nil {
		seelog.Error(err)
		session.Rollback()
		return 0, err
	}
	//3.保存用户添加的出行者到联系人中
	_, err = session.Insert(contacts)
	if err != nil {
		seelog.Error(err)
		session.Rollback()
		return 0, err
	}
	//4.保存订单联系人
	var orderLinks []model.OrderLinker
	for _, elem := range *contacts {
		orderLink := model.OrderLinker{
			OrderID:  order.ID,
			Name:     elem.ContactName,
			Mobile:   elem.ContactMobile,
			IDCard:   elem.IDCard,
			Gender:   elem.Sex,
			BirthDay: elem.Birthday,
		}
		orderLinks = append(orderLinks, orderLink)
	}
	for _, elem := range *outers {
		orderLink := model.OrderLinker{
			OrderID:  order.ID,
			Name:     elem.ContactName,
			Mobile:   elem.ContactMobile,
			IDCard:   elem.IDCard,
			Gender:   elem.Sex,
			BirthDay: elem.Birthday,
		}
		orderLinks = append(orderLinks, orderLink)
	}
	if len(orderLinks) < 1 {
		return 0, errors.New("缺少订单联系人")
	}
	_, err = session.Insert(orderLinks)
	if err != nil {
		seelog.Error(err)
		session.Rollback()
	}
	//提交事务
	session.Commit()
	session.Close()
	//返回订单保存结果
	return val, err
}

//Deprecated
func saveOrderLinkers(session xorm.Session, orderID int64, eleMap map[string]interface{}) (int64, error) {
	idcardStr := eleMap["idcard"].(string)
	//从身份证获取性别
	genderStr := idcardStr[len(idcardStr)-2:len(idcardStr)-1]
	if strings.Contains("24680", genderStr) {
		genderStr = "男"
	} else {
		genderStr = "女"
	}
	//从身份证获取生日
	birthStr := idcardStr[6:14]
	birthDay, err := time.Parse("20060102", birthStr)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	orderLinker := model.OrderLinker{
		OrderID:  orderID,
		Name:     eleMap["name"].(string),
		Mobile:   eleMap["mobile"].(string),
		IDCard:   idcardStr,
		Gender:   genderStr,
		BirthDay: birthDay,
	}
	return session.Insert(&orderLinker)
}

//获取最后提交的一条订单(为了得到订单号)
func GetLastOrder(order *model.Order) (bool, error) {
	return conf.DBEngine.Desc("order_no").Get(order)
}

//根据id获取订单详情
func GetOrderByID(order *model.Order) error {
	orders := []model.Order{}
	err := conf.DBEngine.Where(" id = ? ", order.ID).Find(&orders)
	if err != nil {
		seelog.Error(err)
		return err
	}
	if len(orders) < 1 {
		seelog.Error("没找到订单")
		return errors.New("没找到订单")
	}
	*order = orders[0]
	return nil
}

//根据id获取订单详情
func UpdateOrder(order *model.Order) error {
	_,err := conf.DBEngine.Where(" id = ? ", order.ID).Update(order)
	if err != nil {
		seelog.Error(err)
		return err
	}
	return nil
}

//根据订单id获取订单商品
func GetOrderGoodsByOrderID(orderID int64) (*model.OrderGoods, error) {
	ogs := []model.OrderGoods{}
	og := &model.OrderGoods{}
	err := conf.DBEngine.Where(" order_id = ? ", orderID).Find(&ogs)
	if err != nil {
		seelog.Error(err)
		return og, err
	}
	if len(ogs) < 1 {
		seelog.Error("没有找到订单商品")
		return og, errors.New("没有找到订单商品")
	}

	og = &ogs[0]

	return og, nil
}

//根据订单id获取订单联系人
func GetOrderLinkersByOrderID(orderID int64) ([]model.OrderLinker, error) {
	ols := []model.OrderLinker{}
	err := conf.DBEngine.Where("  dt_orders_id = ? ", orderID).Find(&ols)
	if err != nil {
		seelog.Error(err)
		return ols, err
	}
	return ols, nil
}
