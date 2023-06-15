package main

import "time"

type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	CreationDate time.Time
	Upvotes      int
	Downvotes    int
}
