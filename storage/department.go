package storage

import (
	"fmt"
	"hr-app-back/model"
)

func DepartmentRead(query map[string]string) (departments []model.Department, err error) {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting with database", err)
		return nil, err
	}

	fields := "departmentID, name, description"

	if len(query["fields"]) > 0 {
		fields = query["fields"]
	}

	stmt := session.Select(fields).From("department")

	for k, v := range query {
		if k == "fields" {
			continue
		}
		stmt.Where(fmt.Sprintf("%s = ?", k), v)
	}

	_, err = stmt.Load(&departments)

	return
}
