package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

var (
	gormDB      *gorm.DB
	errDbLoad   error
	errMockLoad error
	db          *sql.DB
	testData    *TestSuite
	postgresDb  *DB
)

func init() {
	connectionString := "host=localhost port=5432 user=recro dbname=recro_test password=recro sslmode=disable"
	gormDB, errDbLoad = gorm.Open("postgres", connectionString)
	postgresDb = &DB{gormDB}
	postgresDb.AutoMigrate(&User{})
}

// mocks the database connection
// TestGormConnection tests gorm connection
func TestGormConnection(t *testing.T) {
	db, mock, errMockLoad := sqlmock.New()
	if errMockLoad != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errMockLoad)
	}
	defer db.Close()
	gormDB, errDbLoad = gorm.Open("postgres", db)
	if errDbLoad != nil {
		t.Errorf("Database not created for testing.")
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
