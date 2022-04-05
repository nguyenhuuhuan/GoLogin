package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
)

type Roles struct {
	gorm.Model
	//RoleId   uint32  `json:"roleId,omitempty" gorm:"column:role_id;private_key;primaryKey" bson:"_id,omitempty" validate:"max=40"`
	RoleName string  `json:"roleName" gorm:"column:role_name;not null;unique" bson:"role_name,omitempty" validate:"required,max=255"`
	Status   string  `json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" `
	Remark   *string `json:"remark,omitempty" gorm:"column:remark" bson:"remark,omitempty" validate:"max=255"`
	CreateBy string  `json:"created_by,omitempty" gorm:"column:created_by" bson:"created_by,omitempty"`
}

func (r *Roles) CreateRole(db *gorm.DB) (*Roles, error) {
	var err error
	err = db.Debug().Create(&r).Error
	if err != nil {
		return &Roles{}, err
	}
	return r, nil
}
func SantizeRole(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}
func (r *Roles) Prepare(userId uint32, db *gorm.DB) {
	user := User{}
	user1, err := user.FindUserById(userId, db)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	r.RoleName = SantizeRole(r.RoleName)
	r.CreateBy = SantizeRole(user1.Username)
	r.Status = "active"
}
func (r *Roles) FindRoleByRoleName(db *gorm.DB, roleName string) (*Roles, error) {
	var err error
	err = db.Debug().Model(&Roles{}).Where("role_name = ?", roleName).Take(&r).Error
	if err != nil {
		return &Roles{}, err
	}
	return r, err
}
