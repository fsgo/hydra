// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package internal

type XServer interface {
	Is(conn XConn) (is bool, err error)
	Head() Head
	Listener() XListener
}

func NewServer(h Head, listener XListener) XServer {
	h.HeaderLen().MustValid()

	return &Server{
		head:     h,
		listener: listener,
	}
}

type Server struct {
	listener XListener
	head     Head
}

func (s *Server) Is(conn XConn) (is bool, err error) {
	hl := s.head.HeaderLen()
	// 先用少量字符做明确的非判断
	header, errHeader := conn.Header(hl[0])
	if errHeader != nil {
		return false, errHeader
	}

	if s.head.MustNot(header) {
		return false, nil
	}

	// 再验证更多字符
	header, errHeader = conn.Header(hl[1])
	if errHeader != nil {
		return false, errHeader
	}

	if s.head.Is(header) {
		return true, nil
	}
	return false, nil
}

func (s *Server) Listener() XListener {
	return s.listener
}

func (s *Server) Head() Head {
	return s.head
}

var _ XServer = &Server{}
