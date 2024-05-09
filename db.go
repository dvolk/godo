package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatal(err)
	}

	sql := `CREATE TABLE IF NOT EXISTS items (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT,
                finished BOOLEAN,
                datetime_added DATETIME
    );`

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InsertTodo(db *sql.DB, name string) {
	stmt, err := db.Prepare("INSERT INTO items (name, finished, datetime_added) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	datetime := time.Now()
	_, err = stmt.Exec(name, false, datetime)
	if err != nil {
		panic(err)
	}
}

func ToggleTodo(db *sql.DB, id string) {
	stmt, err := db.Prepare("UPDATE items set finished = not finished where id = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
}

func DeleteTodo(db *sql.DB, id string) {
	stmt, err := db.Prepare("DELETE FROM items where id = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
}

func AllTodos(db *sql.DB) []Todo {
	rows, err := db.Query("SELECT id, name, finished, datetime_added FROM items")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Name, &todo.Finished, &todo.DatetimeAdded)
		if err != nil {
			panic(err)
		}
		todos = append(todos, todo)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return todos
}
