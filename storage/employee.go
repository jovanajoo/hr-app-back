package storage

import (
	"fmt"
	"hr-app-back/model"
)

func EmployeeInsert(employee *model.Employee) error {

	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting:", err)
		return err
	}

	_, err = session.InsertInto("employee").Columns("employeeID", "organizationID", "firstName", "lastName", "departmentID", "positionID", "phoneNumber", "email", "password", "birthDate", "gender", "address", "employedDate", "remainingVacationDays", "isAdmin").Record(employee).Exec()
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return err
	}

	return nil

}

func EmployeeRead(query map[string]string) (employees []model.Employee, err error) {

	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting database", err)
		return nil, err
	}

	fields := "employeeID, organizationID, firstName, lastName, departmentID, positionID, phoneNumber, email, password, birthDate, gender, address, employedDate, remainingVacationDays, isAdmin, picture"

	if len(query["fields"]) > 0 {
		fields = query["fields"]
	}
	stmt := session.Select(fields).From("employee")

	for k, v := range query {
		if k == "fields" {
			continue
		}

		stmt.Where(fmt.Sprintf("%s = ?", k), v)
	}

	_, err = stmt.Load(&employees)

	return
}

func EmployesGetByOrgId(employee []model.Employee, organizationId int) ([]model.Employee, error) {

	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting database", err)
		return nil, err
	}

	_, err = session.Select("employeeID, organizationID, firstName, lastName, departmentID, positionID, phoneNumber, email, birthDate, gender, address, employedDate, remainingVacationDays, isAdmin").From("employee").Where("organizationID = ?", organizationId).Load(&employee)
	if err != nil {
		fmt.Println("Error getting data:", err)
		return nil, err
	}

	return employee, nil
}

func EmployeeUpdate(employee *model.Employee) error {

	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting database", err)
		return err
	}

	_, err = session.Update("employee").Set("firstName", employee.FirstName).Set("lastName", employee.LastName).Set("departmentID", employee.DepartmentID).Set("positionID", employee.PositionID).Set("phoneNumber", employee.PhoneNumber).Set("email", employee.Email).Set("password", employee.Password).Set("birthDate", employee.BirthDate).Set("gender", employee.Gender).Set("address", employee.Address).Set("employedDate", employee.EmployedDate).Set("RemainingVacationDays", employee.RemainingVacationDays).Where("employeeID = ?", employee.Id).Exec()
	if err != nil {
		return err
	}

	return nil
}

func EmployeeDelete(id int) error {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting database", err)
		return err
	}

	_, err = session.DeleteFrom("employee").Where("employeeID = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}
