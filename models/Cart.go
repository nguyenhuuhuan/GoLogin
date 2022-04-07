package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

var maps = make(map[uint]*CartDTO)

func AddBeverageToCart(db *gorm.DB, cartDTO CartDTO) (*CartDTO, error) {
	item, exist := maps[cartDTO.ID]
	if !exist {
		maps[cartDTO.ID] = &cartDTO
		return maps[cartDTO.ID], nil
	} else {
		b := &Beverage{}
		err := db.Debug().Model(&Beverage{}).Where("id = ?", cartDTO.ID).Take(&b).Error
		if err != nil {
			return &CartDTO{}, err
		}
		if item.Amount >= b.Amount {
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
func RemoveItemCart(id uint) {
	delete(maps, id)
}
func GetAllCart() []*CartDTO {
	var items []*CartDTO
	for _, value := range maps {
		items = append(items, value)
	}
	return items
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
func CreateCarts() ([]*OrderDetail, error) {
	var listOrderDetail []*OrderDetail
	for _, ord := range maps {
		orderDetail := &OrderDetail{}
		orderDetail.NameBev = ord.Name
		orderDetail.ColdHot = ord.HotCold
		orderDetail.Size = ord.Size
		orderDetail.Topping = ord.Topping
		orderDetail.TotalPrice = ord.Total
		listOrderDetail = append(listOrderDetail, orderDetail)
		if listOrderDetail == nil {
			return nil, errors.New("Don't have any item in cart")
		}
	}
	return listOrderDetail, nil
}
