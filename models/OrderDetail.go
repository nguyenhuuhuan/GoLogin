package models

import (
	"github.com/jinzhu/gorm"
	"math/rand"
)

type OrderDetail struct {
	gorm.Model
	Code          string           `gorm:"not null;column:code" json:"code"`
	Size          string           `gorm:"not null;column:size" json:"size"`
	Sugar         string           `gorm:"not null;column:sugar" json:"sugar"`
	ColdHot       string           `gorm:"not null;column:cold_hot" json:"cold_hot"`
	NameBev       string           `gorm:"not null;column:name_bev" json:"name_bev"`
	TotalPrice    float32          `gorm:"not null;column:total_price" json:"total_price"`
	ToppingDetail []*ToppingDetail `gorm:"foreignKey:topping_detail_id;column:topping_detail" json:"topping_detail"`
	Beverage      *Beverage        `gorm:"foreignKey:name;references:name_bev" json:"beverage"`
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = rune(letters[rand.Intn(len(letters))])
	}
	return string(s)
}
func (od *OrderDetail) Prepare() {
	od.ColdHot = Santize(od.ColdHot)
	od.Size = Santize(od.ColdHot)
	od.Sugar = Santize(od.Sugar)
	od.NameBev = Santize(od.NameBev)

}
func (od *OrderDetail) CreateOrderDetail(db *gorm.DB) (*OrderDetail, error) {
	var err error
	err = db.Debug().Create(&od).Error
	if err != nil {
		return &OrderDetail{}, err
	}
	return od, nil
}
func (od *OrderDetail) FindOrderDetailByCode(db *gorm.DB, code string) (*[]OrderDetail, error) {
	var err error
	var listOrderDetail []OrderDetail
	err = db.Debug().Model(&OrderDetail{}).Where("code = ?", code).Take(&listOrderDetail).Error
	if err != nil {
		return nil, err
	}
	return &listOrderDetail, nil
}
