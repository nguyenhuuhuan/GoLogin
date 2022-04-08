package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type CartDTO struct {
	ID            uint
	Name          string
	Amount        uint
	Size          string
	HotCold       string
	Sugar         string
	Price         float32
	ToppingDetail []*ToppingDetail
	Total         float32
}

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
func RemoveCart() {
	for item := range maps {
		delete(maps, item)
	}
}
func GetAllCart() []*CartDTO {
	var items []*CartDTO
	for _, value := range maps {
		items = append(items, value)
	}
	return items
}
func TotalPriceTopping(db *gorm.DB, cartDTO *CartDTO) (float32, error) {
	var totalPriceToppingDetail float32
	var err error
	for _, item := range cartDTO.ToppingDetail {
		topping := &Topping{}
		err = db.Debug().Model(&Topping{}).Where("name_topping = ?", item.NameTop).Take(&topping).Error
		if item.Amount >= topping.Amount {
			return 0, errors.New("Out of range")
		}
		total := topping.Price * float32(item.Amount)
		totalPriceToppingDetail += total
		//item.ToppingDetailID = topping.ID
		topping.Amount = topping.Amount - item.Amount
		item.TotalPrice = totalPriceToppingDetail
	}
	return totalPriceToppingDetail, err
}
func CreateCarts() ([]*OrderDetail, error) {
	var listOrderDetail []*OrderDetail
	for _, ord := range maps {
		orderDetail := &OrderDetail{}
		orderDetail.NameBev = ord.Name
		orderDetail.ColdHot = ord.HotCold
		orderDetail.Size = ord.Size
		orderDetail.ToppingDetail = ord.ToppingDetail
		orderDetail.TotalPrice = ord.Total
		listOrderDetail = append(listOrderDetail, orderDetail)
		if listOrderDetail == nil {
			return nil, errors.New("Don't have any item in cart")
		}
	}
	return listOrderDetail, nil
}
