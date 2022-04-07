package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"web-golang/models"
	"web-golang/response"
	"web-golang/utils/formaterror"
)

func (server *Server) CreateTopping(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	topping := &models.Topping{}
	err = json.Unmarshal(body, &topping)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	toppingCreate, err := topping.CreateTopping(server.DB)
	if err != nil {
		formatError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formatError)
		return
	}
	response.JSON(w, http.StatusOK, toppingCreate)
}
