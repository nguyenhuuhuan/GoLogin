package models

type User_Role struct {
	UserId uint32 `json:"userId,omitempty" gorm:"column:user_id;primary_key" bson:"user_id,omitempty"`
	RoleId uint32 `json:"roleId,omitempty" gorm:"column:role_id;primary_key" bson:"role_id,omitempty"`
}

//err := db.SetupJoinTable(&User{}, "User", &User_Role{})
