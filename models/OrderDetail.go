package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	Code       string    `gorm:"not null;unique;column:code" json:"code"`
	Size       string    `gorm:"not null;column:size" json:"size"`
	Sugar      string    `gorm:"not null;column:sugar" json:"sugar"`
	ColdHot    string    `gorm:"not null;column:cold_hot" json:"cold_hot"`
	TotalPrice float32   `gorm:"not null;column:total_price" json:"total_price"`
	Topping    []Topping `gorm:"" json:"topping"`
}
