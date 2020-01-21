/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/fsgo/fsprotocol/fshead16"

	"github.com/fsgo/hydra"
	"github.com/fsgo/hydra/_examples/server1/repeater"
	"github.com/fsgo/hydra/servers"
	"github.com/fsgo/hydra/servers/fshead16server"
	"github.com/fsgo/hydra/servers/httpserver"
)

func main() {
	s := hydra.NewServer(nil)
	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8090")
	s.SetListenAddr(addr)

	// 注册http 协议
	s.RegisterServer(httpServer())

	// 注册自定义协议
	s.RegisterServer(&repeater.Repeater{})

	// 注册fshead协议server
	s.RegisterServer(fsheadServer())

	err := s.Start()
	log.Fatalln("stopped:", err)
}

func httpServer() servers.Server {
	serveMux := http.NewServeMux()
	server := &httpserver.Server{
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

func fsheadServer() servers.Server {
	fsheadServer := &fshead16server.Server{
		Handler: func(conn net.Conn) {

			var readDeadLine time.Time
			read := func(buf []byte) error {
				conn.SetReadDeadline(readDeadLine)
				_, err := io.ReadFull(conn, buf)
				if err != nil {
					log.Printf("read with error:%v\n", err)
					return err
				}
				return nil
			}

			var writeDeadLine time.Time
			write := func(buf []byte) error {
				conn.SetWriteDeadline(writeDeadLine)
				n, err := conn.Write(buf)
				if err != nil {
					log.Printf("write error:%v\n", err)
					return err
				}
				if n != len(buf) {
					err = fmt.Errorf("wrote only %d bytes, total %d bytes", n, len(buf))
					return err
				}
				return nil
			}

			for {
				readDeadLine = time.Now().Add(1 * time.Second)

				buf := make([]byte, fshead16.Length)
				if err := read(buf); err != nil {
					return
				}
				h, err := fshead16.Load(buf, 0)
				if err != nil {
					log.Printf("parser head error:%v\n", err)
					return
				}
				if h.MetaLen > 0 {
					bufMeta := make([]byte, h.MetaLen)
					if err := read(bufMeta); err != nil {
						return
					}
				}
				bufBody := make([]byte, h.BodyLen)
				if err := read(bufBody); err != nil {
					return
				}

				// 将数据原样写回
				hw := &fshead16.Head{
					BodyLen:    uint32(len(bufBody)),
					ClientName: "server",
				}
				writeDeadLine = time.Now().Add(1 * time.Second)

				if err := write(hw.Bytes()); err != nil {
					return
				}

				if err := write(bufBody); err != nil {
					return
				}
			}
		},
	}
	return fsheadServer
}
