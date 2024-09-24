package models

type Task struct {
	ID         int      `json:"id"`
	Opened     int64    `json:"opened"`
	Closed     int64    `json:"closed"`
	AuthorID   int      `json:"author_id"`
	AssignedID int      `json:"assigned_id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Labels     []string `json:"labels"`
}
