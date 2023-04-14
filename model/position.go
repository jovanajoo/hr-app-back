package model

type Position struct {
	Id          string `db:"positionID" form:"position_id"`
	Name        string `db:"name" form:"name" json:"name"`
	Description string `db:"description" form:"description" json:"description"`
}
