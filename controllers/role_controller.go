package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"web-golang/models"
	"web-golang/response"
	"web-golang/utils/formaterror"
)

func (server *Server) createRole(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	role := models.Roles{}
	err = json.Unmarshal(body, &role)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	param1 := r.FormValue("userId")
	userId, err := strconv.ParseUint(param1, 0, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	role.Prepare(uint32(userId), server.DB)

	roleCreate, err := role.CreateRole(server.DB)
	if err != nil {
		formatError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formatError)
		return
	}
	fmt.Println(r.Body)
	w.Header().Set(" Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, roleCreate.ID))
	response.JSON(w, http.StatusCreated, roleCreate)

}
