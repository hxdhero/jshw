package model

import "time"

//线路幻灯片
type TourismLineImages struct {
	ID           int       `xorm:"id" json:"-"`
	LineID       int       `xorm:"dt_dz_tourismline_id"` //线路id
	OriginalPath string    `xorm:"original_path"`        //原图地址
	ThumbPath    string    `xorm:"thumb_path"`           //缩略图地址
	Remark       string    `xorm:"remark" json:"-"`               //图片描述
	AddTime      time.Time `xorm:"add_time" json:"-"`
}

func (t *TourismLineImages) TableName() string {
	return "dt_dz_tourismlineimages"
}
