package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := AllTodos(App.DB)
	if err != nil {
		log.Println(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{"Todos": todos}
	App.Templates["index.html"].Execute(w, data)
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	expectedFields := []string{"todoName"}
	data, err := PostGetData(w, r, expectedFields)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = InsertTodo(App.DB, data["todoName"])
	if err != nil {
		log.Println(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	expectedFields := []string{"id"}
	data, err := PostGetData(w, r, expectedFields)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = ToggleTodo(App.DB, data["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	expectedFields := []string{"id"}
	data, err := PostGetData(w, r, expectedFields)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = DeleteTodo(App.DB, data["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AllTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := AllTodos(App.DB)
	if err != nil {
		log.Println(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	out, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		http.Error(w, "serialization error", http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
