package models

type ToppingDetail struct {
	ToppingDetailID uint       `gorm:"not null;column:topping_detail_id;primaryKey" json:"topping_detail_id"`
	Code            string     `gorm:"not null;column:code" json:"code"`
	NameTop         string     `gorm:"not null;column:name_top" json:"name_top"`
	Amount          uint       `gorm:"not null;column:amount" json:"amount"`
	TotalPrice      float32    `gorm:"not null;column:total_price" json:"total_price"`
	Topping         []*Topping `gorm:"not null;column:topping;foreignKey:id;References:name_top" json:"topping"`
}
