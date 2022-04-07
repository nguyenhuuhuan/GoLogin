package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"web-golang/auth"
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
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	query := r.URL.Query()
	roleNames, present := query["roleName"]
	if !present || len(roleNames) == 0 {
		response.ERROR(w, http.StatusBadRequest, nil)
		return
	}
	for _, roleName := range roleNames {
		role := &models.Roles{}
		role, err := role.FindRoleByRoleName(server.DB, roleName)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, err)
		}
		user.Roles = append(user.Roles, role)
		//models.AssignRolesToUser(server.DB, *role)
	}
	user.Prepare("register")
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
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreate.ID))
	response.JSON(w, http.StatusCreated, userCreate)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	userId, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenId, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	if tokenId != uint(userId) {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	updateUser, err := user.UpdateUser(uint(userId), server.DB)
	if err != nil {
		//formatError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updateUser)
}
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserById(uint32(uid), server.DB)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	response.JSON(w, http.StatusOK, userGotten)
}
func (server *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.FindAllUser(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, users)
}
