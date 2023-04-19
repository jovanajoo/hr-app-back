package api

import (
	"errors"
	"fmt"
	"hr-app-back/model"
	"hr-app-back/storage"
	"log"
	"net/smtp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func getEmployeeIdByContext(c *gin.Context, id int) (*model.Employee, error) {

	filter := map[string]string{"employeeID": fmt.Sprintf("%d", id)}
	employee, err := storage.EmployeeRead(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee: %v", err)
	}
	if len(employee) == 0 {
		return nil, fmt.Errorf("employe id is not found")
	}

	return &employee[0], nil
}

func getLeaveRequest(c *gin.Context, leaveID int) (*model.Leave, error) {
	filter := map[string]string{"leaveID": fmt.Sprintf("%d", leaveID)}
	leave, err := storage.LeaveRead(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve leave request: %v", err)
	}
	if len(leave) == 0 {
		return nil, fmt.Errorf("leave request not found")
	}

	return &leave[0], nil
}

func getEmpAndOrgByContext(c *gin.Context) (*[]model.Employee, error) {
	orgID, ok := c.Get("organizationID")
	if !ok {
		return nil, errors.New("organization ID is not found in the context")
	}

	filter := map[string]string{"organizationID": fmt.Sprintf("%d", orgID)}
	empByOrg, err := storage.EmployeeRead(filter)
	if err != nil {
		return nil, errors.New("failed to retrieve employee")
	}

	if len(empByOrg) == 0 {
		return nil, errors.New("employees not found")
	}
	return &empByOrg, nil
}

func isAdmin(c *gin.Context) bool {
	isAdmin, ok := c.Get("isAdmin")
	if !ok {
		return false
	}
	return isAdmin.(bool)
}

// func isAuthorized(orgID int, id int) bool {
// 	var employee model.Employee
// 	if employee.OrganizationId != orgID && employee.Id != id {
// 		return false
// 	}
// 	return true
// }

func authEmployee(email string, password string) (*model.Employee, error) {

	employeeFromDB, err := storage.EmployeeRead(map[string]string{"email": email})
	if err != nil {
		return nil, errors.New("invalid email")
	}
	if len(employeeFromDB) == 0 {
		return nil, errors.New("employee is not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(employeeFromDB[0].Password), []byte(password))
	if err != nil {
		return nil, errors.New("bad Email or Password")
	}
	return &employeeFromDB[0], nil
}

func sendPasswordEmail(emp *model.Employee, password string) {
	auth := smtp.PlainAuth("", "milan.tepic221@gmail.com", "qbasglhtdcbqfpkb", "smtp.gmail.com")

	to := []string{"bulatovicj07@gmail.com"}

	msg := []byte("To: bulatovicj07@gmail.com\r\n" +
		"Subject: Why aren’t you using Mailtrap yet?\r\n" +
		"\r\n" +
		"Name:" + emp.FirstName + "\n" +
		"Last Name:" + emp.LastName + "\n" +
		"Email:" + emp.Email + "\n" +
		"Address:" + emp.Address + "\n" +
		"\r\n" +
		"YOUR NEW PASSWORD " + password +
		"Here’s the space for our great sales pitch\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, "milan.tepic221@gmail.com", to, msg)

	if err != nil {
		log.Fatal(err)
	}
}
