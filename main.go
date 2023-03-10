package main

import (
	"fmt"
	"magic/pkg/cronjob"
	"magic/pkg/handler"
	_ "net/http/pprof"

	"magic/app"
	"magic/app/server"
)

func main() {

	go cronjob.Demo()

	fmt.Println("=========================== Application is starting ================== ")
	app.UseBizServer(server.Config{
		Port: 8080,
	}, &handler.DeployHandler{}).Run()

}
