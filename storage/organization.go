package storage

import (
	"errors"
	"fmt"
	"hr-app-back/model"
	"time"
)

func OrganizationRead(query map[string]string) (organizations []model.Organization, err error) {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting with database", err)
		return nil, err
	}

	fields := "organizationID, name, organizationName, organizationEmail, phoneNumber, address, numberOfEmployee, createTime, password"

	if len(query["fields"]) > 0 {
		fields = query["fields"]
	}

	stmt := session.Select(fields).From("organization")

	for k, v := range query {
		if k == "fields" {
			continue
		}
		stmt.Where(fmt.Sprintf("%s = ?", k), v)
	}

	_, err = stmt.Load(&organizations)

	return
}

func OrganizationInsert(organization *model.Organization) error {

	var employee model.Employee
	organization.CreateTime = time.Now()

	session, err := ConnectionToDB()
	if err != nil {
		panic(err)
	}

	/* 	hashPassword, err := utility.HashPassword(organization.Password)
	   	if err != nil {
	   		fmt.Println("Failed to hash password")
	   		return err
	   	}

	   	organization.Password = hashPassword */

	result, err := session.InsertInto("organization").Columns("name", "organizationName", "organizationEmail", "phoneNumber", "address", "numberOfEmployee", "createTime", "password").Record(organization).Exec()
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return err
	}

	employee.OrganizationId = int(id)
	employee.Email = organization.OrganizationEmail
	employee.Password = organization.Password
	employee.IsAdmin = true

	_, err = session.InsertInto("employee").Columns("organizationID", "firstName", "lastName", "phoneNumber", "email", "password", "birthDate", "gender", "address", "employedDate", "remainingVacationDays", "isAdmin").Record(employee).Exec()
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return err
	}

	return nil

}

func OrganizationUpdate(organization *model.Organization) error {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting database", err)
		return err
	}

	_, err = session.Update("organization").Set("firstName", organization.FirstName).Set("lastName", organization.LastName).Set("organizationName", organization.OrganizationName).Set("organizationEmail", organization.OrganizationEmail).Set("phoneNumber", organization.PhoneNumber).Set("address", organization.Address).Set("numberOfEmployee", organization.NumberOfEmployee).Set("CreateTime", organization.CreateTime).Where("organizationID = ?", organization.Id).Exec()
	if err != nil {
		return err
	}
	return nil
}

func OrganizationDelete(id int) error {

	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting database", err)
		return err
	}

	var temp int
	_, err = session.Select("COUNT(*)").From("organization").Where("organizationID = ?", id).Load(&temp)
	if err != nil {
		return err
	}

	if temp == 0 {
		fmt.Println("Id is not found")
		return errors.New("id is not found")
	}

	_, err = session.DeleteFrom("organization").Where("organizationID = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}

/* func InsertRegister(reg *model.Register) (sql.Result, error) {

	reg.CreatedAt = time.Now()

	sess, err := ConnectionToDB("mariadb")
	if err != nil {
		panic(err)
	}

	result, err := sess.InsertInto("registration_organization").Columns("name", "organizationName", "phoneNumber", "address", "numberOfEmployee", "CreatedAt").Record(&register).Exec()
	if err != nil {
		panic(err)
	}

	return result, nil

} */
