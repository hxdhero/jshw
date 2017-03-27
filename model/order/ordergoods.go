package model

//订单商品
type OrderGoods struct {
	ID int64 `xorm:" 'id' pk autoincr <- "`
	OrderID int64 `xorm:"order_id"`
	GoodsTitle string `xorm:"goods_title"`//商品标题
	GoodsPrice float64 `xorm:"goods_price"`//商品价格
	RealPrice float64 `xorm:"real_price"`//实际价格
	Quantity int `xorm:"quantity"`//订购数量
	Point int `xorm:"point"`//积分
	TourismLineID int64 `xorm:"dt_dz_tourismline_id"`//旅游线路id
	TourismLineDateID int64 `xorm:"dt_dz_tourismlinedate_id"`//旅游线路日期id
	Images string `xorm:"images"`//商品图片
}

func (og *OrderGoods) TableName() string{
	return "dt_order_goods"
}
