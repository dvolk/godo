package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return nil, err
	}

	sql := `CREATE TABLE IF NOT EXISTS items (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT,
                finished BOOLEAN,
                datetime_added DATETIME
    );`

	_, err = db.Exec(sql)
	if err != nil {
		return nil, err
	}

	return db, err
}

func InsertTodo(db *sql.DB, name string) error {
	stmt, err := db.Prepare("INSERT INTO items (name, finished, datetime_added) VALUES (?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	datetime := time.Now()
	_, err = stmt.Exec(name, false, datetime)
	if err != nil {
		return err
	}
	return nil
}

func ToggleTodo(db *sql.DB, id string) error {
	stmt, err := db.Prepare("UPDATE items set finished = not finished where id = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTodo(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM items where id = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func AllTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT id, name, finished, datetime_added FROM items")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	todos := []Todo{}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Name, &todo.Finished, &todo.DatetimeAdded)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	err = rows.Err()
	if err != nil {
		return nil, nil
	}
	return todos, nil
}
