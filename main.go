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

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

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
	funcMap := template.FuncMap{
		"urlEncode": urlEncode,
		"icon":      icon,
	}
	todos := map[string]bool{}

	addTodo := func(todoName string) error {
		if _, ok := todos[todoName]; ok {
			return errors.New("todoName already exists")
		}
		todos[todoName] = false
		return nil
	}

	toggleTodo := func(todoName string) error {
		if _, ok := todos[todoName]; !ok {
			return errors.New(fmt.Sprintf("can't find that todoName %s", todoName))
		}
		todos[todoName] = !todos[todoName]
		return nil
	}

	deleteTodo := func(todoName string) error {
		if _, ok := todos[todoName]; !ok {
			return errors.New("can't find that todoName")
		}
		delete(todos, todoName)
		return nil
	}

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			index_tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("index.html")
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			if r.Method != "GET" {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			data := map[string]interface{}{
				"Todos": todos,
			}
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
			expectedFields := []string{"todoName"}
			data, err := PostGetData(w, r, expectedFields)
			if err != nil {
				log.Println(err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			err = toggleTodo(data["todoName"])
			if err != nil {
				log.Println(err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		})

	http.HandleFunc("/todo/delete",
		func(w http.ResponseWriter, r *http.Request) {
			expectedFields := []string{"todoName"}
			data, err := PostGetData(w, r, expectedFields)
			if err != nil {
				log.Println(err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			err = deleteTodo(data["todoName"])
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
