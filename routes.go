package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	todos := allTodos()
	data := map[string]interface{}{"Todos": todos}
	index_tmpl.Execute(w, data)
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	expectedFields := []string{"todoName"}
	data, err := PostGetData(w, r, expectedFields)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = addTodo(data["todoName"])
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	expectedFields := []string{"id"}
	data, err := PostGetData(w, r, expectedFields)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = toggleTodo(data["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	expectedFields := []string{"id"}
	data, err := PostGetData(w, r, expectedFields)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = deleteTodo(data["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AllTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	todos := allTodos()
	out, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
