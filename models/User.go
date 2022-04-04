package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	//UserId   uint32   `gorm:"column:user_id;private_key;auto_increment" json:"user_id,omitempty"`
	Username string   `gorm:"size:255;not null;unique" json:"username"`
	Email    string   `gorm:"size:100;not null;unique" json:"email"`
	Status   string   `gorm:"size:100;not null;unique" json:"status"`
	Password string   `gorm:"size:100;not null;unique" json:"password"`
	Roles    []*Roles `gorm:"many2many:user_role;" json:"roles,omitempty" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func (u *User) BeforeSave() error {
	hashPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}
func Santize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}
func (u *User) Prepare(action string) {
	switch strings.ToLower(action) {
	case "login":
		{
			u.Username = Santize(u.Username)
			u.Email = Santize(u.Email)
			u.Status = "Online"
		}
	case "register":
		{
			u.Username = Santize(u.Username)
			u.Email = Santize(u.Email)
			u.Status = "Offline"
		}
	}

	//for _, roleName := range roleNames {
	//	role := &Roles{}
	//	role, err := role.FindRoleByRoleName(db, roleName)
	//	if err != nil {
	//		fmt.Errorf(err.Error())
	//	}
	//	u.Roles = append(u.Roles, role)
	//}
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Required username")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			if u.Username == "" {
				return errors.New("Required Email or Username")
			}
		} else {
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid Email")
			}
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}

}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	//err1 := u.BeforeSave()
	//if err1 != nil {
	//	log.Fatal(err1)
	//}
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) FindAllUser(db *gorm.DB) (*[]User, error) {
	var err error
	var users []User
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

func (u *User) FindUserById(Id uint32, db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("user_id = ?", Id).Take(&u).Error
	if err != nil {
		return &User{}, errors.New("User not found")
	}
	return u, nil
}
func (u *User) FindUserByUsername(userName string, db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("username = ?", userName).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, err
}

func (u *User) UpdateUser(userId uint, db *gorm.DB) (*User, error) {
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("user_id = ?", userId).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password": u.Password,
			"username": u.Username,
			"email":    u.Email,
			"status":   u.Status,
			//"update_by": u.
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	err = db.Debug().Model(&User{}).Where("user_id = ?", userId).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) deleteUser(userId uint32, db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&User{}).Where("user_id = ?", userId).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//for NO-ORM
//func AssignRolesToUser(db *gorm.DB, user *User, roles []Roles) (err error) {
//	if roles == nil {
//		return nil
//	}
//	for _, role := range roles {
//		err := db.Debug().Create(&User_Role{
//			UserId: user.UserId,
//			RoleId: role.RoleId,
//		}).Error
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
