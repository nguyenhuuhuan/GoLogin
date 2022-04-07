package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Beverage struct {
	gorm.Model
	Name         string  `gorm:"size:255;not null;unique;column:name" json:"name"`
	Amount       uint    `gorm:"not null;column:amount" json:"amount"`
	Price        float32 `gorm:"not null;column:price" json:"price"`
	BeverageType string  `gorm:"not null;column:beverage_type" json:"beverage_type"`
}
type CartDTO struct {
	ID     uint
	Name   string
	Amount uint
	Price  float32
	Total  float32
}
type BeverageType string

const (
	COFFEE   BeverageType = "Coffee"
	SMOOTHIE BeverageType = "Smoothie"
	TEA      BeverageType = "Tea"
	JUICE    BeverageType = "Juice"
)

func (b *Beverage) Prepare() {
	b.Name = Santize(b.Name)
	b.Amount = uint(b.Amount)
}
func (b *Beverage) Validate() error {
	if b.Name == "" {
		return errors.New("Required name")
	}
	if &b.Amount == nil {
		return errors.New("Required amount")
	}
	if &b.Price == nil {
		return errors.New("Required price")
	}
	return nil
}
func (b *Beverage) SaveBeverage(db *gorm.DB) (*Beverage, error) {
	err := db.Debug().Create(&b).Error
	if err != nil {
		return &Beverage{}, err
	}
	return b, nil
}

func (b *Beverage) FindAllBeverage(db *gorm.DB) (*[]Beverage, error) {
	var beverages []Beverage
	err := db.Debug().Model(&Beverage{}).Limit(100).Find(beverages).Error
	if err != nil {
		return &[]Beverage{}, err
	}
	return &beverages, nil
}
func (b *Beverage) FindBeverageById(db *gorm.DB, beverageId uint) (*Beverage, error) {
	var err error
	err = db.Debug().Model(&Beverage{}).Where("id = ?", beverageId).Take(&b).Error
	if err != nil {
		return &Beverage{}, err
	}
	return b, nil
}
func (b *Beverage) FindAllBeverageByType(db *gorm.DB, beverageType string) (*[]Beverage, error) {
	var beverages []Beverage
	err := db.Debug().Model(&Beverage{}).Where("beverage_type = ?", beverageType).Find(&beverages).Error
	if err != nil {
		return &[]Beverage{}, err
	}
	return &beverages, nil
}

var maps = make(map[uint]*CartDTO)

func (b *Beverage) AddBeverageToCart(db *gorm.DB, cartDTO CartDTO) (*CartDTO, error) {
	item, exist := maps[cartDTO.ID]
	if !exist {
		maps[cartDTO.ID] = &cartDTO
		return maps[cartDTO.ID], nil
	} else {
		fmt.Print("hello")
		err := db.Debug().Model(&Beverage{}).Where("id = ?", cartDTO.ID).Take(&b).Error
		if err != nil {
			return &CartDTO{}, err
		}
		if item.Amount > b.Amount {
			return &CartDTO{}, errors.New("Out of range")
		} else {
			fmt.Println("a", item.Amount)
			item.Amount += 1
			item.Total = float32(item.Amount) * item.Price
			return item, nil
		}
	}
	return &CartDTO{}, nil
}

func (b *Beverage) GetAllCart() interface{} {
	var items []interface{}
	for _, value := range maps {
		items = append(items, value)
	}
	return items
}
