package main

import (
	"hr-app-back/api"
)

func main() {

	/* 	_, err := storage.ConnectionToDB()
	   	if err != nil {
	   		panic(err)
	   	} */

	api.SetUpRouters().Run()

}
