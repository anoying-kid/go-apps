package models

import "time"

type Post struct {
    ID        int64     `json:"id"`
    Title     string    `json:"title"`
    Body      string    `json:"body"`
    AuthorID  int64     `json:"author_id"`
    Author    *User     `json:"author,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}