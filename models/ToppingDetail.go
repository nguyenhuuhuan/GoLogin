package models

import "github.com/jinzhu/gorm"

type ToppingDetail struct {
	gorm.Model
	ToppingDetailID uint    `gorm:"not null;column:topping_detail_id" json:"topping_detail_id"`
	Code            string  `gorm:"not null;column:code" json:"code"`
	NameTop         string  `gorm:"not null;column:name_top" json:"name_top"`
	Amount          uint    `gorm:"not null;column:amount" json:"amount"`
	TotalPrice      float32 `gorm:"not null;column:total_price" json:"total_price"`
}
