package controllers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"web-golang/auth"
	"web-golang/models"
	"web-golang/response"
	"web-golang/utils/formaterror"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
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
	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	token, err := server.SignIn(user.Username, user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	response.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(username, email, password string) (string, error) {
	var err error

	user := models.User{}
	var emailOrUsername string
	if email == "" {
		emailOrUsername = username
		err = server.DB.Debug().Model(models.User{}).Where("username = ?", emailOrUsername).Take(&user).Error
	} else {
		emailOrUsername = email
		err = server.DB.Debug().Model(models.User{}).Where("email = ?", emailOrUsername).Take(&user).Error
	}
	if err != nil {
		return "", err
	}
	err = models.CheckPasswordHash(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.UserId)
}
