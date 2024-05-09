package main

import (
	"html/template"
	"log"
	"net/http"
)

var index_tmpl *template.Template
var allTodos func() []Todo
var addTodo func(todoName string) error
var toggleTodo func(todoName string) error
var deleteTodo func(todoName string) error

func main() {
	db := InitDB()

	funcMap := template.FuncMap{
		"icon": Icon,
	}
	var err error
	index_tmpl, err = template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	addTodo = func(todoName string) error {
		InsertTodo(db, todoName)
		return nil
	}

	toggleTodo = func(todoName string) error {
		ToggleTodo(db, todoName)
		return nil
	}

	deleteTodo = func(todoName string) error {
		DeleteTodo(db, todoName)
		return nil
	}

	allTodos = func() []Todo {
		return AllTodos(db)
	}

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/todo/add", AddTodoHandler)
	http.HandleFunc("/todo/toggle", ToggleTodoHandler)
	http.HandleFunc("/todo/delete", DeleteTodoHandler)
	http.HandleFunc("/todo/all", AllTodosHandler)

	serveAddr := ":18080"
	log.Println("trying to listen on", serveAddr)
	log.Fatal(http.ListenAndServe(serveAddr, nil))
}
