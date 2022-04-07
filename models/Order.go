package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	OrderDetail []*OrderDetail `gorm:"foreignKey:code;references:code_bill;column:order_detail" json:"order_detail"`
	TotalBill   float32        `gorm:"column:total_bill;not null" json:"total_bill"`
	CodeBill    string         `gorm:"column:code_bill;not null" json:"code_bill"`
}
