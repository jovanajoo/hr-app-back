package storage

import (
	"fmt"
	"hr-app-back/model"
)

func PositionRead(query map[string]string) (positions []model.Position, err error) {
	session, err := ConnectionToDB()
	if err != nil {
		fmt.Println("Error in connecting with database", err)
		return nil, err
	}

	fields := "positionID, name, description"

	if len(query["fields"]) > 0 {
		fields = query["fields"]
	}

	stmt := session.Select(fields).From("position")

	for k, v := range query {
		if k == "fields" {
			continue
		}
		stmt.Where(fmt.Sprintf("%s = ?", k), v)
	}

	_, err = stmt.Load(&positions)

	return
}
