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
		username := r.FormValue("username")
		password := r.FormValue("password")
		defer db.Close()
		if len(username) > 0 && len(password) > 0 {
			rows, _ := db.Query("SELECT * FROM LOGINS;")
			for rows.Next() {
				rows.Scan(&username, &password)
				fmt.Printf("%s: %s\n", username, password)
			}
		}
		html.ExecuteTemplate(w, "main.html", nil)
	})
	http.ListenAndServe(":8000", nil)
}
