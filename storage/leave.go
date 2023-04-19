package storage

import (
	"fmt"
	"hr-app-back/model"
)

func LeaveRead(query map[string]string) (leave []model.Leave, err error) {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting:", err)
		return nil, err
	}

	fields := "leaveID, categoryName, employeeID, organizationID, startDate, endDate, total, description, status"

	if len(query[fields]) > 0 {
		fields = query[fields]
	}

	stmt := session.Select(fields).From("employee_leave")

	for k, v := range query {
		if k == fields {
			continue
		}

		stmt.Where(fmt.Sprintf("%s = ? ", k), v)
	}

	_, err = stmt.Load(&leave)

	return
}

func LeaveCreate(leave *model.Leave) (int, error) {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting:", err)
		return 0, err
	}

	result, err := session.InsertInto("employee_leave").Columns("leaveID", "categoryName", "employeeID", "organizationID", "startDate", "endDate", "total", "description", "status").Record(leave).Exec()
	if err != nil {
		fmt.Println("Error in inserting data:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func LeaveUpdate(leave *model.Leave) error {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting:", err)
		return err
	}

	_, err = session.Update("employee_leave").Set("categoryName", leave.CategoryName).Set("employeeID", leave.EmployeeId).Set("organizationID", leave.OrganizationId).Set("startDate", leave.StartDate).Set("total", leave.Total).Set("description", leave.Description).Set("status", leave.Status).Where("leaveID = ?", leave.LeaveID).Exec()
	if err != nil {
		return err
	}

	return nil
}

func LeaveUpdateStatus(leave *model.Leave) error {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting:", err)
		return err
	}

	_, err = session.Update("employee_leave").Set("status", leave.Status).Where("leaveID = ?", leave.LeaveID).Exec()
	if err != nil {
		return err
	}

	return err
}
