package main

import "time"

type Todo struct {
	ID            int
	Name          string
	Finished      bool
	DatetimeAdded time.Time
}
