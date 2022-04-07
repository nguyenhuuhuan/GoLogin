package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("nickname already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("email already taken")
	}
	if strings.Contains(err, "title") {
		return errors.New("title already taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect Password")
	}
	if strings.Contains(err, "role_name") {
		return errors.New("roleName already taken")
	}
	if strings.Contains(err, "name_topping") {
		return errors.New("Topping name already taken")
	}
	return errors.New("incorrect details")
}
