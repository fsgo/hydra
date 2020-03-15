/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/21
 */

package xnet

import (
	"github.com/fsgo/hydra/protocol"
)

type Server interface {
	Protocol() protocol.Protocol
	Is(conn XConn) (is bool, err error)
	Listener() XListener
	Serve() error
	Close() error
}

func newServer(p protocol.Protocol, listenerSize int) Server {
	p.HeaderLen().MustValid()

	return &server{
		protocol: p,
		listener: NewListener(listenerSize),
	}
}

type server struct {
	listener XListener
	protocol protocol.Protocol
}

func (s *server) Is(conn XConn) (is bool, err error) {
	hl := s.protocol.HeaderLen()
	// 先用少量字符做明确的非判断
	header, errHeader := conn.Header(hl[0])
	if errHeader != nil {
		return false, errHeader
	}

	if s.protocol.MustNot(header) {
		return false, nil
	}

	// 在验证更多字符
	header, errHeader = conn.Header(hl[1])
	if errHeader != nil {
		return false, errHeader
	}

	if s.protocol.Is(header) {
		return true, nil
	}
	return false, nil
}

func (s *server) Close() error {
	return s.protocol.Close()
}

func (s *server) Serve() error {
	return s.protocol.Serve(s.listener)
}

func (s *server) Listener() XListener {
	return s.listener
}

func (s *server) Protocol() protocol.Protocol {
	return s.protocol
}

var _ Server = &server{}
