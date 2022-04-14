package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	OrderDetail   []*OrderDetail `gorm:"foreignKey:code;references:code_bill;column:order_detail" json:"order_detail"`
	TotalBeverage float32        `gorm:"column:total_beverage;not null" json:"total_beverage"`
	TotalTopping  float32        `gorm:"column:total_topping;not null" json:"total_topping"`
	TotalBill     float32        `gorm:"column:total_bill;not null" json:"total_bill"`
	CodeBill      string         `gorm:"column:code_bill;not null" json:"code_bill"`
}

func (o *Order) SaveOrder(db *gorm.DB) (*Order, error) {
	var err error
	err = db.Debug().Create(&o).Error
	if err != nil {
		return nil, err
	}
	return o, nil
}
