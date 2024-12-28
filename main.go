package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var tmpl *template.Template
var db *sql.DB

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")

}

func initDB() {
	var err error
	// Initialize the db variable
	db, err = sql.Open("mysql", "root:root@(127.0.0.1:5426)/testdb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	gRouter := mux.NewRouter()

	//Setup MySQL
	initDB()
	defer db.Close()

	http.ListenAndServe(":4000", gRouter)

}
