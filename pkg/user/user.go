package user

import "time"

type User struct {
	Id        string
	Password  string
	Email     string
	VelogName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
