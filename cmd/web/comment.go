package main

import "time"

type Comment struct {
	ID           int
	UserID       int
	PostID       int
	Content      string
	CreationDate time.Time
}
