/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package main

import (
	"log"
	"net"
	"net/http"

	fshead2 "github.com/fsgo/fshead"

	"github.com/fsgo/mpserver"
	"github.com/fsgo/mpserver/_examples/server1/protocol/repeater"
	"github.com/fsgo/mpserver/protocol"
	"github.com/fsgo/mpserver/protocol/fshead"
	httpProtocol "github.com/fsgo/mpserver/protocol/http"
)

func main() {
	s := mpserver.NewServer(nil)
	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8090")
	s.SetListenAddr(addr)

	// 注册http 协议
	s.RegisterProtocol(httpServer())

	// 注册自定义协议
	s.RegisterProtocol(&repeater.Protocol{})

	// 注册fshead协议server
	s.RegisterProtocol(fsheadServer())

	err := s.Start()
	log.Fatalln("stopped:", err)
}

func httpServer() protocol.Protocol {
	serveMux := http.NewServeMux()
	server := &httpProtocol.Protocol{
		Server: &http.Server{
			Handler: serveMux,
		},
	}
	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("你好:"))
		writer.Write([]byte(request.RequestURI))
	})
	return server
}

func fsheadServer() protocol.Protocol {
	fsheadServer := &fshead.Protocol{
		Handler: func(h *fshead2.FsHead, metaRead fshead.Read, bodyRead fshead.Read) []byte {
			log.Println("clientName=", h.ClientName)
			log.Println("logID=", h.LogID)
			log.Println("userID=", h.UserID)
			body, err := bodyRead()
			if err != nil {
				return []byte(err.Error())
			}
			return body
		},
	}
	return fsheadServer
}
