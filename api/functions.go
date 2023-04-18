package api

import (
	"errors"
	"fmt"
	"hr-app-back/model"
	"hr-app-back/storage"

	"github.com/gin-gonic/gin"
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

func getEmployeeByContext(c *gin.Context) (*model.Employee, int, error) {
	orgID, ok := c.Get("organizationID")
	if !ok {
		return nil, 0, errors.New("organization ID is not found in the context")
	}

	email, ok := c.Get("email")
	if !ok {
		return nil, 0, errors.New("email not found in context")
	}

	filter := map[string]string{"organizationID": fmt.Sprintf("%d", orgID), "email": email.(string)}
	employeeFromDB, err := storage.EmployeeRead(filter)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve employee")
	}

	if len(employeeFromDB) == 0 {
		return nil, 0, errors.New("employee not found")
	}
	return &employeeFromDB[0], orgID.(int), nil
}

func isAdmin(c *gin.Context) bool {
	isAdmin, ok := c.Get("isAdmin")
	if !ok {
		return false
	}
	return isAdmin.(bool)
}

/* func isAuthorized(c *gin.Context, orgID int, id int) bool {
	employee, err := getEmployeeByContext(c)
	if err != nil {
		return false
	}

	return employee.OrganizationId == orgID && employee.Id == id
}
*/
