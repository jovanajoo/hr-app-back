package model

import "time"

type Leave struct {
	LeaveID        int       `db:"leaveID"`
	CategoryName   string    `db:"categoryName" form:"category_name" json:"category_name"`
	EmployeeId     int       `db:"employeeID" validate:"required"`
	OrganizationId int       `db:"organizationID" validate:"required"`
	StartDate      time.Time `db:"startDate" form:"start_date" json:"start_date" validate:"required"`
	EndDate        time.Time `db:"endDate" form:"end_date" json:"end_date" validate:"required,gtfield=StartDate"`
	Total          int       `db:"total" form:"total" json:"total,omitempty"`
	Description    string    `db:"description" form:"description" json:"description,omitempty"`
	Status         string    `db:"status" json:"status" validate:"required,oneof=pending approved rejected"`
}
