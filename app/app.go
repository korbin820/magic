package app

import (
	"service-tag/app/server"
	"service-tag/app/server/biz"
)

func UseWebServer(config server.Config, handler server.Handler) server.Server {
	s := &server.WebServer{}

	s.Config = config

	s.Handler = handler

	return s
}

func UseBizServer(config server.Config, handler biz.BizHandler) server.Server {
	s := &biz.BizServer{}

	s.Config = config

	s.Handler = handler

	return s
}
