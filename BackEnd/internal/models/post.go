package models

import "time"

type Post struct {
	Id         int       `json:"id,omitempty"`
	UserId     string    `json:"user_id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Subject    string    `json:"subject,omitempty"`
	ParentId   int       `json:"parent_id,omitempty"`
	Categories []int     `json:"categories,omitempty"`
	Comments   []Post    `json:"comments,omitempty"`
}

type Mark struct {
	PostId int    `json:"post_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Mark   bool   `json:"mark,omitempty"`
}

type PostAndMarks struct {
	Post
	Likes      int        `json:"likes"`
	Dislikes   int        `json:"dislikes"`
	Categories []Category `json:"categories,omitempty"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
