package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func NewApplication(db *sql.DB) *Application {
	var err error
	funcMap := template.FuncMap{
		"icon": Icon,
	}
	index_tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	templateMap := map[string]*template.Template{
		"index.html": index_tmpl}

	app := Application{
		DB:        db,
		Templates: templateMap,
		FuncMap:   funcMap,
	}

	return &app
}

var App *Application

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatal("couldn't initialize database")
	}
	App = NewApplication(db)

	http.HandleFunc("/", Chain(IndexHandler, Method("GET"), Logging()))
	http.HandleFunc("/todo/add", Chain(AddTodoHandler, Method("POST"), Logging()))
	http.HandleFunc("/todo/toggle", Chain(ToggleTodoHandler, Method("POST"), Logging()))
	http.HandleFunc("/todo/delete", Chain(DeleteTodoHandler, Method("POST"), Logging()))
	http.HandleFunc("/todo/all", Chain(AllTodosHandler, Method("GET"), Logging()))

	serveAddr := ":18080"
	log.Println("trying to listen on", serveAddr)
	log.Fatal(http.ListenAndServe(serveAddr, nil))
}
