package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"
)

type Application struct {
	DB        *sql.DB
	Templates map[string]*template.Template
	FuncMap   template.FuncMap
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Todo struct {
	ID            int
	Name          string
	Finished      bool
	DatetimeAdded time.Time
}
