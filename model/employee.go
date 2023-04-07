package model

import "time"

type Employee struct {
	Id                    int       `db:"employeeID" form:"id" json:"id"`
	OrganizationId        int       `db:"organizationID" form:"organization_id" json:"organization_id"`
	FirstName             string    `db:"firstName" form:"first_name" json:"first_name"`
	LastName              string    `db:"lastName" form:"last_name" json:"last_name"`
	DepartmentID          string    `db:"departmentID" form:"department_id" json:"department_id"`
	PositionID            string    `db:"positionID" form:"position_id" json:"position_id"`
	PhoneNumber           string    `db:"phoneNumber" form:"phone_number" json:"phone_number"`
	Email                 string    `db:"email" form:"email" json:"email"`
	Password              string    `db:"password"  form:"password" json:"password"`
	BirthDate             time.Time `db:"birthDate" form:"birth_date" json:"birth_date"`
	Gender                string    `db:"gender" form:"gender" json:"gender"`
	Address               string    `db:"address" form:"address" json:"address"`
	EmployedDate          time.Time `db:"employedDate" form:"employed_date" json:"employed_date"`
	RemainingVacationDays int       `db:"remainingVacationDays" form:"remaining_vacation_days" json:"remaining_vacation_days"`
	IsAdmin               bool      `db:"isAdmin"`
}
