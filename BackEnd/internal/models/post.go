package models

import (
	"fmt"
	"time"
)

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
	UserLogin  string `json:"user_login,omitempty"`
	Likes      int    `json:"likes,omitempty"`
	Dislikes   int    `json:"dislikes,omitempty"`
	Categories string `json:"categories,omitempty"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CommentsAndMarks struct {
	Post
	UserLogin  string             `json:"user_login,omitempty"`
	Likes      int                `json:"likes,omitempty"`
	Dislikes   int                `json:"dislikes,omitempty"`
	Categories string             `json:"categories,omitempty"`
	Children   []CommentsAndMarks `json:"children"`
}

func (t *CommentsAndMarks) AddNestedChild(arr []CommentsAndMarks) {
	remember := arr[0].Id

	for i, comment := range arr {
		fmt.Print(remember)
		if i > 0 && comment.ParentId == remember {
			t.Children = append(t.Children, comment)
			fmt.Print("added: ", remember, comment.Id)
		} else {
			if len(arr) != i+1 {
				if remember < arr[i+1].ParentId {
					remember = arr[i+1].ParentId
				}
			}
			if len(arr) == i {
				t.Children = append(t.Children, comment)
				fmt.Print("added: ", remember, comment.Id)
			}
		}

		fmt.Println()
	}

	//fmt.Printf("%+v\n", &arr)
}

//// If this child is one level below the current node, just add it here for now
//if newEntry.Level == t.Level+1 {
//t.Children = append(t.Children, &newEntry)
//} else {
//// Loop through the children and see if it fits anywhere
//for _, child := range t.Children {
//if newEntry.Left > child.Left && newEntry.Right < child.Right {
//child.AddNestedChild(newEntry)
//break
//}
//}
