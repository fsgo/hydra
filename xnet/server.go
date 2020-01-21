/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/21
 */

package xnet

import (
	"github.com/fsgo/hydra/servers"
)

type Server interface {
	Server() servers.Server
	Listener() Listener
	Serve() error
	Close() error
}

func newServer(ss servers.Server, lsSize int) Server {
	return &server{
		ServerVal:   ss,
		ListenerVal: NewListener(lsSize),
	}
}

type server struct {
	ListenerVal Listener
	ServerVal   servers.Server
}

func (s *server) Close() error {
	return s.ServerVal.Close()
}

func (s *server) Serve() error {
	return s.ServerVal.Serve(s.ListenerVal)
}

func (s *server) Listener() Listener {
	return s.ListenerVal
}

func (s *server) Server() servers.Server {
	return s.ServerVal
}

var _ Server = &server{}
