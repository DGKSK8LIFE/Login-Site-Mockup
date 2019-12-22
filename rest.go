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
	loginSite  *template.Template
	createSite *template.Template
)

func init() {
	loginSite = template.Must(template.ParseGlob("login.html"))
	createSite = template.Must(template.ParseGlob("create.html"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loginSite.ExecuteTemplate(w, "main.html", nil)
	})
	http.HandleFunc("/create.html", func(w http.ResponseWriter, r *http.Request) {
		createSite.ExecuteTemplate(w, "create.html", nil)
	})
	http.HandleFunc("/login", userAuth)
	http.HandleFunc("/create", createAccount)
	http.ListenAndServe(":8000", nil)
}

func userAuth(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	db, err := sql.Open("sqlite3", "stored_logins.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if row := rowExists(username, password, db); row == true {
		fmt.Fprint(w, "<h1 style='text-align: center;'>Welcome!</h1>")
	} else if row == false {
		// deny access
		loginSite.ExecuteTemplate(w, "login.html", nil)
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	db, err := sql.Open("sqlite3", "stored_logins.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO logins (username, password) VALUES ('%s', '%s');", username, password)
	db.Exec(query)
}

func rowExists(username, password string, db *sql.DB) bool {
	var exists bool
	query := fmt.Sprintf("SELECT * FROM ACCOUNTS WHERE username='%s' AND password='%s'", username, password)
	if err := db.QueryRow(query).Scan(&username, &password); err != nil && err != sql.ErrNoRows {
		log.Fatal("database error, we're fucked")
	} else if err == sql.ErrNoRows {
		exists = false
	} else {
		exists = true
	}
	return exists
}
