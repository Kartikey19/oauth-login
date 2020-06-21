package postgres

import (
	"database/sql/driver"
	"encoding/json"

	_ "github.com/lib/pq"
)

// JSONB defining custom type for map[string]interafce{}
type JSONB map[string]interface{}

// Value overiding default function for updating a field value
func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// Scan overriding default function for fetching a field's value from database
func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

// CREATE INDEX index_users_on_name ON users USING gin(to_tsvector('simple', name));

// User represents app user
type User struct {
	ID    int64  `gorm:"primary_key" json:"id"`
	Name  string `gorm:"not null type:tsvector" json:"name"`
	Email string `gorm:"unique" json:"email"`
	Phone string `gorm:"index" json:"phone"`
	Meta  JSONB  `sql:"type:jsonb" json:"-"`
}

// SearchUser object will returned when searching users
type SearchUser struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CheckValidSource checks if the source is valid
func (db *DB) CheckValidSource(source string) bool {
	switch source {
	case "twitter",
		"facebook",
		"github":
		return true
	}
	return false
}

// CreateUser creates a new user object in database
func (db *DB) CreateUser(user *User) int64 {

	if db.NewRecord(user) {
		db.Create(&user)
		return user.ID
	}
	return -1
}

// CheckUserExists checks is user is already in system for a given email id
func (db *DB) CheckUserExists(email string) *User {
	var user User
	db.Where("email = ?", email).Find(&user)
	return &user
}

// UpdateUserMeta updates the user meta data
func (db *DB) UpdateUserMeta(userID int64, meta map[string]interface{}) error {
	user := User{ID: userID}
	err := db.Model(&user).Update("meta", meta).Error
	return err
}

// GetAllUsers fetches all users from database
func (db *DB) GetAllUsers() []*User {
	var users []*User
	db.Model(&User{}).Find(&users)
	return users
}

// GetUserByID searches user by id
func (db *DB) GetUserByID(id int64) *User {
	var user User
	db.Where(&User{ID: id}).First(&user)
	return &user
}

// SearchUserByName searches users by name
func (db *DB) SearchUserByName(term string) ([]SearchUser, error) {
	var results []SearchUser
	err := db.Table("users").Where("name @@ to_tsquery(?)", term).Find(&results).Error
	return results, err
}
