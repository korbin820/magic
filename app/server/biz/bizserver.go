package biz

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"time"

	"service-tag/app/server"
	util "service-tag/utils"

	jsoniter "github.com/json-iterator/go"
)

type BizServer struct {
	Config  server.Config
	Handler BizHandler
}

type BizHandler interface {
	DoGet(context context.Context, path string, params map[string][]string) interface{}
	DoPost(context context.Context, path string, body []byte) interface{}
}

func (s *BizServer) Run() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), &echo{server: s}); err != nil {
	//	fmt.Println(fmt.Sprintf("listen err:%s", err.Error()))
	//	panic(err)
	//}

	addr := ":8080"
	if s.Config.Port != 0 {
		addr = ":" + fmt.Sprint(s.Config.Port)
	}
	svr := http.Server{
		Addr:           addr,
		Handler:        &echo{server: s},
		IdleTimeout:    4 * 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	svr.SetKeepAlivesEnabled(true)

	fmt.Println("listening .....:", s.Config.Port)

	svr.ListenAndServe()
}

type echo struct {
	server *BizServer
}

type Response struct {
	Error   int
	Message string
	Data    interface{}
}

func prepareQueryString(route *url.URL) (string, map[string][]string) {

	if route == nil {
		panic(NewError("route is nil"))
	}

	rawQuery, err := url.QueryUnescape(route.RawQuery)

	if err != nil {
		panic(NewError("rawQuery parse err"))
	}

	var params = make(map[string][]string)

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		panic(NewError(err.Error()))
	}
	for key, value := range values {
		params[key] = value
	}

	return strings.ToLower(route.Path), params
}

func (h echo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.Method == "POST" {
		if req.Body == nil {
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		}
	}

	defer func() {

		if r := recover(); r != nil {
			switch r.(type) {
			case runtime.Error:
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "app::server::go runtime error:", r, "\n", string(buf[:n]))
				var response = Response{
					Error:   999,
					Message: "runtime error",
					Data:    nil,
				}
				res := util.StructToJson(response)
				w.Write([]byte(res))
			case Error: // 业务异常
				bizError := r.(Error)
				var response = Response{
					Error:   1,
					Message: bizError.Error(),
					Data:    nil,
				}
				res := util.StructToJson(response)
				w.Write([]byte(res))
			default: // 非运行时错误
				fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "serve error:", r)
				var response = Response{
					Error:   999,
					Message: "server error",
					Data:    nil,
				}
				res := util.StructToJson(response)
				w.Write([]byte(res))
			}
		}
	}()

	ctx := context.WithValue(context.Background(), "headers", req.Header)

	path, para := prepareQueryString(req.URL)

	if path == "/" {
		w.Write([]byte("ok"))
		return
	}

	if path == "/healthcheck" {
		w.Write([]byte("ok"))
		return
	}

	if ok, _ := regexp.MatchString("/html/", req.URL.String()); ok {
		http.StripPrefix("/html/", http.FileServer(http.Dir("./html/"))).ServeHTTP(w, req)
		return
	}

	if req.Method == "GET" {
		//path, params := prepareQueryString(req.URL)
		fmt.Println("时间:" + time.Now().Format(time.RFC3339) + ",Get request ====> path:" + path + ",para:" + util.StructToJson(para))
		res := h.server.Handler.DoGet(ctx, path, para)

		var response = Response{
			Error:   0,
			Message: "",
			Data:    res,
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(util.StructToJson(response)))
	} else if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		fmt.Println("时间:" + time.Now().Format(time.RFC3339) + ",Post request ====> path:" + path + ",para:" + string(body))
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		//path, para := prepareQueryString(req.URL)

		//如果带有参数auth_module，auth_prop以及参数格式为auth_prop="repoName"&auth_module="monitor/monitorapi"&auth_module="monitor/logdata"
		if len(para) > 0 {
			err := checkModule(para, body)

			if err != nil {
				panic(NewError(err.Error()))
			}
		}

		res := h.server.Handler.DoPost(ctx, path, body)

		var response = Response{
			Error:   0,
			Message: "",
			Data:    res,
		}

		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(util.StructToJson(response)))
	}

}

func checkModule(para map[string][]string, body []byte) error {
	prop := ""
	modules := make([]string, 0)

	if len(para["auth_prop"]) > 0 {
		prop = para["auth_prop"][0]
	} else {
		return nil
	}

	if len(para["auth_module"]) > 0 {
		modules = para["auth_module"]
	} else {
		return nil
	}

	moduleName := jsoniter.Get(body, prop).ToString()

	fmt.Println("prop:" + prop + ",moduleName:" + moduleName + ",modules:" + util.StructToJson(modules))

	checked := false
	for _, v := range modules {
		if strings.HasPrefix(moduleName, v) {
			checked = true
			break
		}
	}

	if !checked {
		return errors.New("无此模块权限")
	}

	return nil
}
