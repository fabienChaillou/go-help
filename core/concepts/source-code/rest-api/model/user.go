package model

type User struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" binding:"required" db:"name"`
}
