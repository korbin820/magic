package main

import (
	"fmt"
	"magic/pkg/cronjob"
	"magic/pkg/handler"
	"magic/pkg/zookeeper"
	_ "net/http/pprof"

	"magic/app"
	"magic/app/server"
)

func main() {

	go cronjob.Demo()

	fmt.Println(zookeeper.NewManager().Get("/test"))

	fmt.Println("=========================== Application is starting ================== ")
	app.UseBizServer(server.Config{
		Port: 8080,
	}, &handler.DeployHandler{}).Run()

}
