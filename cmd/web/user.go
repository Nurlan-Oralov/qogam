package main

import "time"

type User struct {
	ID               int
	Username         string
	Email            string
	Password         string
	RegistrationDate time.Time
}
