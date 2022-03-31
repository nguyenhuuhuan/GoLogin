package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Roles struct {
	RoleId   uint32    `json:"roleId,omitempty" gorm:"column:role_id;private_key" bson:"_id,omitempty" dynamodbav:"roleId,omitempty" firestore:"roleId,omitempty" validate:"max=40"`
	RoleName string    `json:"roleName,omitempty" gorm:"column:roleName" bson:"roleName,omitempty" dynamodbav:"roleName,omitempty" firestore:"roleName,omitempty" validate:"required,max=255"`
	Status   string    `json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty" match:"equal" validate:"required,max=1,code"`
	Remark   *string   `json:"remark,omitempty" gorm:"column:remark" bson:"remark,omitempty" dynamodbav:"remark,omitempty" firestore:"remark,omitempty" validate:"max=255"`
	CreateBy string    `json:"createdBy,omitempty" gorm:"column:createdBy" bson:"createdBy,omitempty" dynamodbav:"createdBy,omitempty" firestore:"createdBy,omitempty"`
	CreateAt time.Time `json:"createdAt,omitempty" gorm:"column:createdAt" bson:"createdAt,omitempty" dynamodbav:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	UpdateBy *string   `json:"updatedBy,omitempty" gorm:"column:updatedBy" bson:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty" firestore:"updatedBy,omitempty"`
	UpdateAt time.Time `json:"updatedAt,omitempty" gorm:"column:updatedAt" bson:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
	//Privileges []string   `json:"privileges,omitempty" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
}
type RoleModule struct {
	RoleId      string `json:"roleId,omitempty" gorm:"column:roleId" bson:"roleId,omitempty" dynamodbav:"roleId,omitempty" firestore:"roleId,omitempty" validate:"required"`
	ModuleId    string `json:"moduleId,omitempty" gorm:"column:moduleId" bson:"moduleId,omitempty" dynamodbav:"moduleId,omitempty" firestore:"moduleId,omitempty" validate:"required"`
	Permissions int32  `json:"permissions,omitempty" gorm:"column:permissions" bson:"permissions,omitempty" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty" validate:"required"`
}

func (r *Roles) CreateRole(db *gorm.DB) (*Roles, error) {
	var err error
	err = db.Debug().Create(&Roles{}).Error
	if err != nil {
		return &Roles{}, err
	}
	return r, nil

}
func (r *Roles) Prepare(userId uint32, db *gorm.DB) {
	//r.RoleId += 1
	user := User{}
	user1, err := user.FindUserById(userId, db)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	r.RoleName = "ROLE_USER"
	r.CreateAt = time.Now()
	r.UpdateAt = time.Now()
	r.CreateBy = user1.Username
	r.UpdateBy = nil
	r.Status = "active"
}
