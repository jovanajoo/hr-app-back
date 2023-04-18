package model

type Category struct {
	Name        string `db:"name"`
	Description string `db:"description"`
}
