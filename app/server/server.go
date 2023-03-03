package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"runtime"
	"time"
)

type Config struct {
	Port    int64
	TimeOut int64
}

type Server interface {
	Run()
}

type Handler interface {
	Handle(uri string, data []byte) ([]byte, error)
}

type WebServer struct {
	Config  Config
	Handler Handler
}

func (s *WebServer) Run() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("listening .....:", s.Config.Port)

	svr := http.Server{
		Addr:           fmt.Sprintf(":%d", s.Config.Port),
		Handler:        &echo{server: s},
		IdleTimeout:    4 * 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	svr.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "connection")
			// 新的连接，计数加1
			//s.wg.Add(1)
		case http.StateActive:
			// 有新的请求，从idle conn pool中移除
			//s.mu.Lock()
			//delete(s.conns, conn.LocalAddr().String())
			//s.mu.Unlock()
		case http.StateIdle:
			//select {
			//case <-s.quit:
			//	// 如果要关闭了，直接Close，否则加入idle conn pool中。
			//	conn.Close()
			//default:
			//	s.mu.Lock()
			//	s.conns[conn.LocalAddr().String()] = conn
			//	s.mu.Unlock()
			//}
		case http.StateHijacked, http.StateClosed:
			// conn已经closed了，计数减一
			//s.wg.Done()
			fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "close")
		}
	}

	svr.ListenAndServe()
}

type echo struct {
	server *WebServer
}

func (h echo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {

		if r := recover(); r != nil {
			switch r.(type) {
			case runtime.Error:
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "app::server::go runtime error:", r, "\n", string(buf[:n]))
				w.Write([]byte(string(buf[:n])))
			default: // 非运行时错误
				fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "serve error:", r)
				w.Write([]byte("error"))
			}
		}
	}()

	url := path.Clean(req.URL.Path)

	if req.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	if url == "/" {
		w.Write([]byte("ok"))
	} else {

		buf, err := h.server.Handler.Handle(url, body)

		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		w.Write(buf)
	}

}
