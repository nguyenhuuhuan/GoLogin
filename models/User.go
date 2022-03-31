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
	UserId   uint32    `gorm:"private_key;auto_increment" json:"user_id"`
	Username string    `gorm:"size:255;not null;unique" json:"username"`
	Email    string    `gorm:"size:100;not null;unique" json:"email"`
	Password string    `gorm:"size:100;not null;unique" json:"password"`
	CreateAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
	Roles    []Roles   `gorm:"many2many:user_role" json:"roles,omitempty" bson:"roles,omitempty" dynamodbav:"roles,omitempty" firestore:"roles,omitempty"`
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func (u *User) BeforSave() error {
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
func (u *User) Prepare() {
	u.UserId += 1
	u.Username = Santize(u.Username)
	u.Email = Santize(u.Email)
	u.UpdateAt = time.Now()
	u.CreateAt = time.Now()
	//u.Roles = nil
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
	err1 := u.BeforSave()
	if err1 != nil {
		log.Fatal(err1)
	}
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

func (u *User) UpdateUser(userId uint32, db *gorm.DB) (*User, error) {
	err := u.BeforSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("user_id = ?", userId).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"username":  u.Username,
			"email":     u.Email,
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
