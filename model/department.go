package model

type Department struct {
	Id          string `db:"departmentID"`
	Name        string `db:"name" form:"name" json:"name"`
	Description string `db:"description" form:"description" json:"description"`
}
