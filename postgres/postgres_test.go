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

func init() {
	connectionString := "host=localhost port=5432 user=recro dbname=recro_test password=recro sslmode=disable"
	gormDB, errDbLoad = gorm.Open("postgres", connectionString)
	db = &DB{gormDB}
	db.AutoMigrate(&User{})
}

// TestGormConnection tests gorm connection
func TestGormConnection(t *testing.T) {
	if errDbLoad != nil {
		t.Errorf("Database not created for testing.")
	}
}
