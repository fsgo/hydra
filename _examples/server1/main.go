// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/28

package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/fsgo/hydra"
	"github.com/fsgo/hydra/_examples/server1/repeater"
	"github.com/fsgo/hydra/protocols"
)

func main() {
	s := &hydra.Hydra{}

	s.BindServer(protocols.HTTP(), httpServer())

	s.BindServer(&repeater.Protocol{}, &repeater.Server{})

	ln, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("now hydra server Listen:", ln.Addr().Network(), ln.Addr().String())

	hydra.WaitShutdown(s, 10*time.Second)

	err = s.Serve(ln)

	log.Fatalln("stopped:", err)
}

func httpServer() *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("你好:"))
		writer.Write([]byte(request.RequestURI))
	})
	return &http.Server{
		Handler: serveMux,
	}
}
