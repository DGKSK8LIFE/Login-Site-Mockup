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
		html.ExecuteTemplate(w, "main.html", nil)
		if len(username) > 0 && len(password) > 0 {
			rows, _ := db.Query("SELECT * FROM LOGINS;")
			for rows.Next() {
				err := rows.Scan(&username, &password)
				if err != nil {
					panic(err)
				} else {
					if r.FormValue("username") == username && r.FormValue("password") == password {
						fmt.Fprint(w, "<h1 style='text-align: center'>welcome home blyat!</h1>")
					} else {
						fmt.Fprintf(w, "<h1 style='text-align: center'>GTFO!</h1>")
					}
				}

			}
		}

	})
	http.ListenAndServe(":8000", nil)
}
