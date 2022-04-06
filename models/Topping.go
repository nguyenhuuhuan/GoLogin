package models

import "github.com/jinzhu/gorm"

type Topping struct {
	gorm.Model
	NameTopping string  `gorm:"size:255;not null;unique;column:name_topping" json:"nameTopping"`
	Amount      uint    `gorm:"not null;column:amount" json:"amount"`
	Price       float32 `gorm:"not null;column:price" json:"price"`
}
