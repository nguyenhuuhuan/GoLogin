package models

type User_Role struct {
	UserId uint32 `json:"userId,omitempty" gorm:"column:userId;primary_key" bson:"_id,omitempty"`
	RoleId uint32 `json:"roleId,omitempty" gorm:"column:roleId;primary_key" bson:"_id,omitempty" dynamodbav:"roleId,omitempty" firestore:"roleId,omitempty"`
}
