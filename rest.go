package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var html *template.Template

func init() {
	html = template.Must(template.ParseGlob("main.html"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db, _ := sql.Open("sqlite3", "./stored_logins.sqlite")
		usernameUserSide := r.FormValue("username")
		passwordUserSide := r.FormValue("password")
		defer db.Close()
		html.ExecuteTemplate(w, "main.html", nil)
		if len(usernameUserSide) > 0 && len(passwordUserSide) > 0 {
			rows, _ := db.Query("SELECT * FROM LOGINS;")
			var username string
			var password string
			for rows.Next() {
				err := rows.Scan(&username, &password)
				if err != nil {
					panic(err)
				} else {
					if usernameUserSide == username && passwordUserSide == password {
						fmt.Fprintln(w, "<h1 style='text-align: center;'>welcome back, cyka blyat!</h1>")
						break
					}
				}
			}

		}
	})
	http.ListenAndServe(":8000", nil)
}
