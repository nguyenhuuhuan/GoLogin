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

	//var roles []models.Roles
	role := models.Roles{}
	query := r.URL.Query()
	roleNames, present := query["roleName"]
	if !present || len(roleNames) == 0 {
		response.ERROR(w, http.StatusBadRequest, nil)
		return
	}
	for _, roleName := range roleNames {
		_, err := role.FindRoleByRoleName(server.DB, roleName)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, err)
		}
		//roles = append(role)
	}
	//err := models.AssignRolesToUser(server.DB, &user, roles)
	user.Prepare()
	err = user.Validate("")
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
