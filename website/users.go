package website

import (
	"fmt"
	"html/template"
	"net/http"
	"recro_demo/postgres"
	"strconv"

	"github.com/go-chi/chi"
)

// UserPageData template response for all users
type UserPageData struct {
	Users []*postgres.User
}

func (web *Website) fetchAllUsers(w http.ResponseWriter, r *http.Request) {
	users := web.Env.DB.GetAllUsers()
	resp := &UserPageData{Users: users}
	tmpl := template.Must(template.ParseFiles("templates/users.html"))
	err := tmpl.Execute(w, resp)
	if err != nil {
		fmt.Fprint(w, `<p>There is some problem fetching all Users!</p>`)
	}
}

func (web *Website) fetchUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userIntID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		fmt.Fprint(w, `<p>There is some problem fetching all Users!</p>`)
	}
	dbUser := web.Env.DB.GetUserByID(userIntID)
	if dbUser.ID == 0 {
		fmt.Fprint(w, `<p>Not a valid user!</p>`)
	}
	t, err := template.New("User").Parse("You are user \"{{ .Name}}\" with email: \"{{ .Email}}\"")
	err = t.Execute(w, dbUser)
	if err != nil {
		fmt.Fprint(w, `<p>There is some problem fetching all Users!</p>`)
	}
}
