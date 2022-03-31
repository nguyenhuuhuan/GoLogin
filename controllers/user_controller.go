package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"web-golang/models"
	"web-golang/response"
	"web-golang/utils/formaterror"
)

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	role := models.Roles{}
	user.Password, err = models.Hash(user.Password)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	//username, err := user.FindUserByUsername(user.Username, server.DB)
	//if err == nil {
	//	response.ERROR(w, http.StatusBadRequest, fmt.Errorf("Username %s is already ", username.Username))
	//}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	role.Prepare(1, server.DB)
	_, err = role.CreateRole(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreate, err := user.SaveUser(server.DB)
	if err != nil {
		formatError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formatError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreate.UserId))
	response.JSON(w, http.StatusCreated, userCreate)

}
