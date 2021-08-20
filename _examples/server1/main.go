// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/28

package main

import (
	"log"
	"net"
	"net/http"

	"github.com/fsgo/hydra"
	"github.com/fsgo/hydra/_examples/server1/repeater"
	"github.com/fsgo/hydra/xhead/xhttp"
)

func main() {
	s := hydra.New(hydra.OptionsDefault)

	// 注册http 协议
	httpLn, _ := s.BindHead(&xhttp.Head{})
	go serveHTTP(httpLn)

	// 注册自定义协议
	rpLn, _ := s.BindHead(&repeater.Head{})

	go serveRepeater(rpLn)

	ln, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("now hydra server Listen:", ln.Addr().Network(), ln.Addr().String())

	err = s.Serve(ln)

	log.Fatalln("stopped:", err)
}

func serveHTTP(ln net.Listener) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("你好:"))
		writer.Write([]byte(request.RequestURI))
	})
	err := http.Serve(ln, serveMux)
	log.Println("http server exit:", err)
}

func serveRepeater(ln net.Listener) {
	rps := &repeater.Server{}
	rps.Serve(ln)
	log.Println()
}
