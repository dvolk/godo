package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
)

func icon(s string) template.HTML {
	return template.HTML(fmt.Sprintf(`<i class="fa fa-fw fa-%s"></i>`, s))
}

func PostGetData(w http.ResponseWriter, r *http.Request, expectedFields []string) (map[string]string, error) {
	// works with both json and web form requests
	out := map[string]string{}
	if r.Method != "POST" {
		return out, errors.New("wrong method")
	}
	if r.Header.Get("Content-Type") == "application/json" {
		// it's a json request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return out, errors.New("can't read body")
		}
		err = json.Unmarshal(body, &out)
		if err != nil {
			return out, errors.New("can't parse body as json")
		}
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return out, errors.New("Can't parse form")
		}
		for key, values := range r.Form {
			fmt.Println(key, values)
			for _, str := range values {
				val, err := url.QueryUnescape(str)
				if err != nil {
					return out, errors.New(fmt.Sprintf("couldn't unescape %s", str))
				}
				out[key] = val
			}
		}
	}

	for _, expectedField := range expectedFields {
		if _, ok := out[expectedField]; !ok {
			return out, errors.New("missing field")
		}
	}
	return out, nil
}

func main() {
	db := InitDB()

	funcMap := template.FuncMap{
		"icon": icon,
	}

	addTodo := func(todoName string) error {
		InsertTodo(db, todoName)
		return nil
	}

	toggleTodo := func(todoName string) error {
		ToggleTodo(db, todoName)
		return nil
	}

	deleteTodo := func(todoName string) error {
		DeleteTodo(db, todoName)
		return nil
	}

	allTodos := func() []Todo {
		return AllTodos(db)
	}

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			// move this into main to avoid remaking template on every page load
			index_tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			if r.Method != "GET" {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			todos := AllTodos(db)
			data := map[string]interface{}{"Todos": todos}
			index_tmpl.Execute(w, data)
		})

	http.HandleFunc("/todo/add",
		func(w http.ResponseWriter, r *http.Request) {
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
		})

	http.HandleFunc("/todo/toggle",
		func(w http.ResponseWriter, r *http.Request) {
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
		})

	http.HandleFunc("/todo/delete",
		func(w http.ResponseWriter, r *http.Request) {
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
		})

	http.HandleFunc("/todo/all",
		func(w http.ResponseWriter, r *http.Request) {
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
		})

	serveAddr := ":18080"
	log.Println("trying to listen on", serveAddr)
	log.Fatal(http.ListenAndServe(serveAddr, nil))
}
