package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	html    *template.Template
	htmlTwo *template.Template
)

func init() {
	html = template.Must(template.ParseGlob("main.html"))
	htmlTwo = template.Must(template.ParseFiles("create.html"))
}

func main() {
	http.HandleFunc("/", showLoginPage)
	http.HandleFunc("/create.html/", showCreateSite)
	http.ListenAndServe(":8000", nil)
}

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./stored_logins.sqlite")
	if err != nil {
		log.Fatal(err)
	}
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
			}
			if usernameUserSide == username && passwordUserSide == password {
				fmt.Fprintln(w, "<h1 style='text-align: center;'>welcome back, cyka blyat!</h1>")
				break
			}
		}
	}
}

func showCreateSite(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./stored_logins.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	usernameClientSide := r.FormValue("username")
	passwordClientSide := r.FormValue("password")
	defer db.Close()
	htmlTwo.Execute(w, "create.html")

	if len(usernameClientSide) > 0 && len(passwordClientSide) > 0 {
		insertQ := fmt.Sprintf("INSERT INTO LOGINS (username, password)\nVALUES ('%s', '%s');", usernameClientSide, passwordClientSide)
		prepQ, err := db.Prepare(insertQ)
		if err != nil {
			log.Fatal(err)
		}
		prepQ.Exec()

		fmt.Fprintf(w, "<h1 style='text-align: center;'>You made an account!</h1>")
	}
}
