package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"web-golang/models"
	"web-golang/response"
)

func (server *Server) createBeverage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	beverage := models.Beverage{}
	err = json.Unmarshal(body, &beverage)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	beverage.Prepare()
	err = beverage.Validate()
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	beverageCreate, err := beverage.SaveBeverage(server.DB)
	if err != nil {
		if strings.Contains(err.Error(), "name") {
			err = errors.New("Name beverage already taken")
		}
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, beverageCreate)
}
func (server *Server) GetAllBeverage(w http.ResponseWriter, r *http.Request) {
	beverage := models.Beverage{}
	beverages, err := beverage.FindAllBeverage(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, beverages)
}
func (server *Server) GetBeveragesByType(w http.ResponseWriter, r *http.Request) {
	beverage := models.Beverage{}
	param := r.FormValue("beverageType")

	beverageTypes, err := beverage.FindAllBeverageByType(server.DB, param)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
	}
	response.JSON(w, http.StatusOK, beverageTypes)
}
