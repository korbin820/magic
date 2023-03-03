package handler

import (
	"context"
	_ "net/http/pprof"
	"strings"
)

type DeployHandler struct {
}

func (handler *DeployHandler) DoGet(ctx context.Context, path string, params map[string][]string) interface{} {

	var response interface{}

	if strings.HasPrefix(path, "/monitor-tag") {
		path = strings.Replace(path, "/monitor-tag", "", 1)
	}

	switch path {

	case "/get/test":

		response = "get test ~"

	}

	return response
}

func (handler *DeployHandler) DoPost(ctx context.Context, path string, data []byte) interface{} {

	var response interface{}

	if strings.HasPrefix(path, "/monitor-tag") {
		path = strings.Replace(path, "/monitor-tag", "", 1)
	}

	switch path {

	/*
	 * method：post
	 * params： null
	 * action：测试接口
	 */
	case "/post/test":

		response = "post test ~"

	}

	return response
}
