package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Topping struct {
	gorm.Model
	NameTopping string  `gorm:"size:255;not null;unique;column:name_topping" json:"nameTopping"`
	Amount      uint    `gorm:"not null;column:amount" json:"amount"`
	Price       float32 `gorm:"not null;column:price" json:"price"`
}

func (t *Topping) CreateTopping(db *gorm.DB) (*Topping, error) {
	var err error
	err = db.Debug().Create(&t).Error
	if err != nil {
		return &Topping{}, err
	}
	return t, err
}
func (t *Topping) UpdateTopping(db *gorm.DB, toppingId uint) (*Topping, error) {
	db = db.Debug().Model(&Topping{}).Where("id = ?", toppingId).Take(&Topping{}).UpdateColumn(
		map[string]interface{}{
			"name_topping": t.NameTopping,
			"amount":       t.Amount,
		},
	)
	if db.Error != nil {
		return &Topping{}, db.Error
	}
	var err error
	err = db.Debug().Model(&Topping{}).Where("id = ?", toppingId).Take(&t).Error
	if err != nil {
		return &Topping{}, err
	}
	return t, nil
}
func (t *Topping) FindToppingByName(db *gorm.DB, nameTop string) (*Topping, error) {
	var err error
	err = db.Debug().Model(&Topping{}).Where("name_topping = ?", nameTop).Take(&t).Error
	if err != nil {
		return &Topping{}, err
	}
	return t, nil
}
func AddToppingToBeverageCart(db *gorm.DB, cartDTO *CartDTO) {
	for _, nameTop := range cartDTO.Topping {
		topping := &Topping{}
		topping, err := topping.FindToppingByName(db, nameTop.NameTopping)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		cartDTO.Topping = append(cartDTO.Topping, topping)
	}
}
func TotalPriceTopping(db *gorm.DB, cartDTO *CartDTO) float32 {
	var totalPriceTopping float32
	var err error
	for _, item := range cartDTO.Topping {
		topping := &Topping{}
		err = db.Debug().Model(&Topping{}).Where("name_topping = ?", item.NameTopping).Take(&topping).Error
		if err != nil {
			fmt.Errorf(err.Error())
		}
		item.ID = topping.ID
		item.Price = topping.Price
		total := topping.Price * float32(item.Amount)
		totalPriceTopping += total
	}
	return totalPriceTopping
}
