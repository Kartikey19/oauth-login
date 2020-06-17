package postgres

import (
	"github.com/jinzhu/gorm"
)

// DB - client to access the postgres database
type DB struct {
	*gorm.DB
}

// EnableVerboseMode enables logging on database query
func (db *DB) EnableVerboseMode() {
	db.LogMode(true)
}

// InitPostgresDB - Initialises and returns a DB client
func InitPostgresDB(ConnectionString string) (*DB, error) {
	db, err := gorm.Open("postgres", ConnectionString)
	// db.LogMode(true)
	if err != nil {
		return nil, err
	}
	postgreDB := DB{db}
	postgreDB.AutoMigrate(&User{})
	return &postgreDB, err
}
