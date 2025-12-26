package task

import "time"

type Task struct {
	Id          int       `json:"id"`
	Taskname    string    `json:"taskname"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Due         time.Time `json:"due"`
	Done        bool      `json:"done"`

	Assignments []*Assignment
}
