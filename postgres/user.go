package postgres

import (
	"encoding/json"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// User represents app user
type User struct {
	ID    int64          `gorm:"primary_key" json:"id"`
	Name  string         `gorm:"not null" json:"name"`
	Email string         `gorm:"unique" json:"email"`
	Phone string         `gorm:"index" json:"phone"`
	Meta  postgres.Jsonb `json:"meta"`
}

func (db *DB) checkValidSource(source string) bool {
	switch source {
	case "twitter",
		"facebook",
		"github":
		return true
	}
	return false
}

func (db *DB) createUser(user *User) int64 {

	if db.NewRecord(user) {
		db.Create(&user)
		return user.ID
	}
	return -1
}

func (db *DB) checkUserExists(email string) *User {
	var user User
	db.Where("email = ?", email).Find(&user)
	return &user
}

func (db *DB) updateUserMeta(userID int64, meta json.RawMessage) error {
	var user User
	err := db.First(&user, userID).Update("meta", meta).Error
	return err
}
