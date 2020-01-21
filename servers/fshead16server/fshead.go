/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/12
 */

package fshead16server

import (
	"net"
	"sync/atomic"

	"github.com/fsgo/fsprotocol/fshead16"

	"github.com/fsgo/hydra/servers"
)

// Server HTTP 协议
type Server struct {
	Head    fshead16.Head
	Handler func(conn net.Conn)
	running int32
}

func (p *Server) HeaderLen() int {
	return p.Head.DiscernLen()
}

func (p *Server) Is(header []byte) bool {
	return p.Head.Is(header)
}

func (p *Server) Name() string {
	return "fshead16"
}

func (p *Server) Serve(l net.Listener) error {
	p.running = 1
	for {
		if atomic.LoadInt32(&p.running) != 1 {
			return nil
		}

		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go p.Handler(conn)
	}
}

func (p *Server) Close() error {
	atomic.StoreInt32(&p.running, 0)
	return nil
}

var _ servers.Server = &Server{}
