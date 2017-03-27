package model

import "time"

//订单
type Order struct {
	ID              int64     `xorm:" 'id' pk autoincr <- "`
	OrderNO         string    `xorm:"order_no"`
	TradeNO         string    `xorm:"trade_no"` //交易号担保支付时使用
	UserID          int     `xorm:"user_id"`  //
	UserName        string    `xorm:"user_name"`
	PaymentID       int64     `xorm:"payment_id"`      //支付方式
	PaymentFee      float64   `xorm:"payment_fee"`     //支付手续费
	PaymentStatus   int       `xorm:"payment_status"`  //支付状态 1.未支付 2.已支付
	PaymentTime     time.Time `xorm:"payment_time"`    //支付时间
	Mobile          string    `xorm:"mobile"`          //手机号
	Message         string    `xorm:"message"`         //描述
	RealAmount      float64   `xorm:"real_amount"`     //实际支付
	OrderAmount     float64   `xorm:"order_amount"`    //订单价格
	Status          int       `xorm:"status"`          //1.订单生成2.订单确认
	AddTime         time.Time `xorm:"add_time"`        //下单时间
	PlayNum         int       `xorm:"playnum"`         //出行人数
	ApplyRefundDate time.Time `xorm:"applyrefunddate"` //申请退款日期
	RefundReason    string    `xorm:"refundreason"`    //退款理由
	BuyerEmail      string    `xorm:"buyer_email"`     //下单人邮箱
	OutDate         time.Time `xorm:"outdate"`         //出行日期
	BackDate        time.Time `xorm:"backdate"`        //返回日期
	RendezvousID    int       `xorm:"rendezvousid"`    //集合地点id dt_dz_point
	RendezvousName  string    `xorm:"rendezvousname"`  //集合地点名
}

func (o *Order) TableName() string {
	return "dt_orders"
}
