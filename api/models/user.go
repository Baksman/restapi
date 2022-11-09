package models

import (
	"errors"
	// "fmt"
	"html"
	"log"
	"strings"
	"time"

	// "github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	// "gorm.io/gorm"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;unique;not null" json:"nickname" `
	Email     string    `gorm:"size:100;unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"size:100;unique;not null" json:"password" validate:"required,min=6,max=60"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := HashPassword(u.Password)

	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	// u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// func (u *User) Validate(actions string) error {
// 	switch strings.ToLower(actions) {
// 	case "update":
// 		if u.Nickname == "" {
// 			return errors.New("required nickname")
// 		}
// 		if u.Password == "" {
// 			return errors.New("required Password")
// 		}
// 		if u.Email == "" {
// 			return errors.New("required Email")
// 		}
// 		if err := checkmail.ValidateFormat(u.Email); err != nil {
// 			return errors.New("invalid Email")
// 		}
// 		return nil

// 	case "login":
// 		if u.Password == "" {
// 			return errors.New("required Password")
// 		}
// 		if u.Email == "" {
// 			return errors.New("required Email")
// 		}
// 		if err := checkmail.ValidateFormat(u.Email); err != nil {
// 			return errors.New("invalid email")
// 		}
// 		return nil

// 	default:
// 		if u.Nickname == "" {
// 			return errors.New("required nickname")
// 		}
// 		if u.Password == "" {
// 			return errors.New("required password")
// 		}
// 		if u.Email == "" {
// 			return errors.New("required email")
// 		}
// 		if err := checkmail.ValidateFormat(u.Email); err != nil {
// 			return errors.New("invalid email")
// 		}
// 		return nil

// 	}
// }

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error

	users := []User{}

	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {

	err := db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error

	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User not found")
	}
	return u, err
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {

	err := u.BeforeSave()

	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"nickname":   u.Nickname,
			"email":      u.Email,
			"updated_at": time.Now(),
		})

	if db.Error != nil {
		return &User{}, db.Error
	}

	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
