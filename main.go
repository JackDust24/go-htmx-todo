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

type Task struct {
	Id   int
	Task string
	Done bool
}

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")

}

func initDB() {
	var err error
	db, err = sql.Open("pgx", "root:root@(127.0.0.1:5432)/crud?parseTime=true")
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

	gRouter.HandleFunc("/", Homepage)

	//Get Tasks
	gRouter.HandleFunc("/tasks", fetchTasks).Methods("GET")

	//Fetch Add Task Form
	gRouter.HandleFunc("/newtaskform", getTaskForm)

	http.ListenAndServe(":4000", gRouter)

}

func Homepage(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "home.html", nil)

}

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	todos, _ := getTasks(db)
	//fmt.Println(todos)

	//If you used "define" to define the template, use the name you gave it here, not the filename
	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func getTaskForm(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "addTaskForm", nil)
}

func getTasks(dbPointer *sql.DB) ([]Task, error) {

	query := "SELECT id, task, done FROM tasks"

	rows, err := dbPointer.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var tasks []Task

	for rows.Next() {
		var todo Task

		rowErr := rows.Scan(&todo.Id, &todo.Task, &todo.Done)

		if rowErr != nil {
			return nil, err
		}

		tasks = append(tasks, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil

}
