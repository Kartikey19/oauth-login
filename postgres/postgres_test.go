package postgres

import (
	"testing"

	"github.com/jinzhu/gorm"
)

var (
	gormDB    *gorm.DB
	errDbLoad error
	db        *DB
	testData  *TestSuite
)

// Gorm db will be initialised in memory rather than a file
func init() {
	connectionString := "host=localhost port=5432 user=recro dbname=recro_test password=recro sslmode=disable"
	gormDB, errDbLoad = gorm.Open("postgres", connectionString)
	db = &DB{gormDB}
	db.AutoMigrate(&User{})

	testData = &TestSuite{}
	u1 := &TestUser{
		Name:   "t1",
		Email:  "t1@gmail.com",
		Phone:  "8283948439",
		Source: "github",
		Other:  "oos",
	}

	u2 := &TestUser{
		Name:   "t1",
		Email:  "t1@gmail.com",
		Phone:  "8283948439",
		Source: "facebook",
		Other:  "oos",
	}
	u3 := &TestUser{
		Name:   "t2",
		Email:  "t2@gmail.com",
		Phone:  "8283948439",
		Source: "github",
		Other:  "oos",
	}

	testData.Users = []*TestUser{u1, u2, u3}
}

// TestGormConnection tests gorm connection
func TestGormConnection(t *testing.T) {
	if errDbLoad != nil {
		t.Errorf("Database not created for testing.")
	}
}
