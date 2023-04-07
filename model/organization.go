package model

import "time"

type Organization struct {
	Id                int       `db:"organizationID" form:"id" json:"id"`
	Name              string    `db:"name"  form:"employee_name" json:"employee_name"`
	OrganizationName  string    `db:"organizationName"  form:"organization_name" json:"organization_name"`
	OrganizationEmail string    `db:"organizationEmail"  form:"organization_email" json:"organization_email"`
	PhoneNumber       int       `db:"phoneNumber" form:"phone_number" json:"phone_number"`
	Address           string    `db:"address" json:"address"`
	NumberOfEmployee  int       `db:"numberOfEmployee"  form:"number_of_employee" json:"number_of_employee"`
	CreateTime        time.Time `db:"createTime"  form:"create_time" json:"create_time"`
	Password          string    `db:"password"  form:"password" json:"password"`
}
