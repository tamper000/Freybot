package models

type User struct {
	ID       int64 `gorm:"primaryKey"`
	Model    string
	Provider string
	Group    string
	Photo    string
	Role     string
	Edit     string
}
