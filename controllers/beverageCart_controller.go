package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"web-golang/models"
	"web-golang/response"
)

var amount *uint

func (server *Server) addBeverageToCart(w http.ResponseWriter, r *http.Request) {
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
	cartDTO.ID = beverageGotten.ID
	cartDTO.Name = beverageGotten.Name
	cartDTO.Price = beverageGotten.Price
	cartDTO.Amount = 1
	beverageGotten.Amount = beverageGotten.Amount - cartDTO.Amount
	cartDTO.Total = float32(cartDTO.Amount) * cartDTO.Price
	createItem, err := beverage.AddBeverageToCart(server.DB, cartDTO)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	response.JSON(w, http.StatusOK, createItem)
}

func (server *Server) GetAllCart(w http.ResponseWriter, r *http.Request) {
	beverage := models.Beverage{}
	carts := beverage.GetAllCart()
	response.JSON(w, http.StatusOK, carts)
}
