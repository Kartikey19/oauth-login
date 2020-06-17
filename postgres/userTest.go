package postgres

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

var (
	gormDB    *gorm.DB
	errDbLoad error
	db        *DB
	testData  *TestSuite
)

type TestUser struct {
	Name   string
	Email  string
	Phone  string
	Other  string
	Source string
}

type TestSuite struct {
	Users []*TestUser
}

func init() {
	connectionString := "host=localhost port=5432 user=recro dbname=recro_test password=recro sslmode=disable"
	gormDB, errDbLoad = gorm.Open("postgres", connectionString)
	db = &DB{gormDB}
	db.AutoMigrate(&User{})
}

// TestUserCreation tests first time user creation
func TestUserCreation(t *testing.T) {
	testData := &TestSuite{}
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
	for _, user := range testData.Users {
		fmt.Println(user)
		// check if email already exists
		dbUser := db.checkUserExists(user.Email)
		var metaMap map[string]interface{}
		if dbUser.ID == -1 {
			// create new user
			dbUser = &User{
				Name:  user.Name,
				Email: user.Email,
				Phone: user.Phone,
			}

			metaMap[user.Source] = user
			e, err := json.Marshal(metaMap)
			if err != nil {
				t.Error("New User: Not marshaling second time json")
				return
			}
			metadata := json.RawMessage(e)
			dbUser.Meta = postgres.Jsonb{RawMessage: metadata}
			id := db.createUser(dbUser)
			if id == -1 {
				t.Error("New User: Valid user not created")
			}
		} else {
			// update the current user
			var jsonData map[string]interface{}
			err := json.Unmarshal(dbUser.Meta.RawMessage, &jsonData)
			if err != nil {
				t.Error("Not unmarshaling old json data")
				return
			}
			if _, ok := jsonData[user.Source]; ok {
				e, err := json.Marshal(user)
				if err != nil {
					t.Error("Repeated User:Not marshal json")
					return
				}
				jsonData[user.Source] = e

			} else {
				jsonData[user.Source] = user

			}
			e, err := json.Marshal(jsonData)
			if err != nil {
				t.Error("Repeated User:Not marshal json")
				return
			}
			metadata := json.RawMessage(e)
			err = db.updateUserMeta(dbUser.ID, metadata)
			if err != nil {
				t.Error("Repeated User:Meta not updated")
				return
			}
		}

	}
}
