package main

import (
	"fmt"
	_ "net/http/pprof"
	"service-tag/pkg/handler"

	"service-tag/app"
	"service-tag/app/server"
)

func main() {

	// go corn.UpdateProjectTags()

	fmt.Println("=========================== Application is starting ================== ")
	app.UseBizServer(server.Config{
		Port: 8080,
	}, &handler.DeployHandler{}).Run()

}
