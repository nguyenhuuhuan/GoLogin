package models

import (
	"github.com/jinzhu/gorm"
	"math/rand"
)

type OrderDetail struct {
	gorm.Model
	Code    string `gorm:"not null;unique;column:code" json:"code"`
	Size    string `gorm:"not null;column:size" json:"size"`
	Sugar   string `gorm:"not null;column:sugar" json:"sugar"`
	ColdHot string `gorm:"not null;column:cold_hot" json:"cold_hot"`
	NameBev string `gorm:"not null;column:name_bev" json:"name_bev"`
	//NameTop    []string   `gorm:"not null;column:name_top;type:text[]" json:"name_top"`
	TotalPrice float32    `gorm:"not null;column:total_price" json:"total_price"`
	Topping    []*Topping `gorm:"foreignKey:name_topping;references:topping" json:"topping"`
	Beverage   *Beverage  `gorm:"foreignKey:name;references:name_bev" json:"beverage"`
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
func (od *OrderDetail) Prepare() {
	od.Code = Santize(RandomString(5))
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
