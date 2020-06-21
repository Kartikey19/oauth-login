package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
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
	u4 := &TestUser{
		Name:   "samarth",
		Email:  "samarth@gmail.com",
		Phone:  "8283948459",
		Source: "github",
		Other:  "oos",
	}

	testData.Users = []*TestUser{u1, u2, u3, u4}
	for _, user := range testData.Users {
		// check if email already exists
		dbUser := postgresDb.CheckUserExists(user.Email)
		if dbUser.ID == 0 {
			// create new user
			createNewUser(user, t)
		} else {
			updateUser(user, t, dbUser)
		}

	}
}

func createUserStruct(user *TestUser) *User {
	dbUser := &User{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
	return dbUser
}

func createNewUser(user *TestUser, t *testing.T) {
	dbUser := createUserStruct(user)
	metaMap := make(map[string]interface{})
	metaMap[user.Source] = user

	dbUser.Meta = metaMap
	id := postgresDb.CreateUser(dbUser)
	if id == -1 {
		t.Error("New User: Valid user not created")
	}
}

func updateUser(user *TestUser, t *testing.T, dbUser *User) {
	// update the current user
	jsonData := dbUser.Meta
	if jsonData == nil {
		jsonData = make(map[string]interface{})
	}
	if _, ok := jsonData[user.Source]; !ok {
		jsonData[user.Source] = structs.Map(&user)
		err := postgresDb.UpdateUserMeta(dbUser.ID, jsonData)
		if err != nil {
			t.Error("Repeated User:Meta not updated")
			return
		}
	}

}

func TestAllUsersFetching(t *testing.T) {
	db, mock, errMockLoad := sqlmock.New()
	if errMockLoad != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errMockLoad)
	}

	gormDB, errDbLoad = gorm.Open("postgres", db)
	if errDbLoad != nil {
		t.Errorf("Database not created for testing.")
	}
	defer gormDB.Close()
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(`SELECT * from "users"`)

}
