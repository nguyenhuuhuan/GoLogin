package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"web-golang/models"
	"web-golang/response"
)

var amount *uint

func (server *Server) addBeverageToCart(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 32, 10)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	beverage := models.Beverage{}
	beverageGotten, err := beverage.FindBeverageById(server.DB, uint(uid))
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	cartDTO := models.CartDTO{}
	err = json.Unmarshal(body, &cartDTO)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cartDTO.ID = beverageGotten.ID
	cartDTO.Name = beverageGotten.Name
	cartDTO.Price = beverageGotten.Price
	cartDTO.Amount = 1
	beverageGotten.Amount = beverageGotten.Amount - cartDTO.Amount
	priceTopping := models.TotalPriceTopping(server.DB, &cartDTO)
	if cartDTO.Size == "M" {
		cartDTO.Total = float32(cartDTO.Amount)*cartDTO.Price + priceTopping + 5000
	} else if cartDTO.Size == "L" {
		cartDTO.Total = float32(cartDTO.Amount)*cartDTO.Price + priceTopping + 7000
	} else {
		cartDTO.Total = float32(cartDTO.Amount)*cartDTO.Price + priceTopping
	}

	createItem, err := models.AddBeverageToCart(server.DB, cartDTO)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	response.JSON(w, http.StatusOK, createItem)
}
func (server *Server) RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	models.RemoveItemCart(uint(uid))
	response.JSON(w, http.StatusOK, "Item was removed successfully")
}
func (server *Server) GetAllCart(w http.ResponseWriter, r *http.Request) {
	carts := models.GetAllCart()
	response.JSON(w, http.StatusOK, carts)
}
func (server *Server) SaveCart(w http.ResponseWriter, r *http.Request) {
	//var listOrderDetail []*models.OrderDetail
	listOrderDetail, err := models.CreateCarts()
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var listOrdDetail []*models.OrderDetail
	code := models.RandomString(5)

	for _, ordDetail := range listOrderDetail {
		ordDetail.Code = code
		createOrdDetail, err := ordDetail.CreateOrderDetail(server.DB)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		listOrdDetail = append(listOrderDetail, createOrdDetail)
	}
	response.JSON(w, http.StatusOK, listOrdDetail)
}
