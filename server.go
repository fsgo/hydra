// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package hydra

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type picker interface {
	Pick(n int) ([]byte, error)
}

func newServer(h Protocol, listener *listener) *server {
	h.DiscernLengths().MustValid()

	return &server{
		head:     h,
		listener: listener,
	}
}

type server struct {
	listener *listener
	head     Protocol
}

func (s *server) Is(conn picker) (is bool, err error) {
	hl := s.head.DiscernLengths()
	// 先用少量字符做明确的非判断
	header, errHeader := conn.Pick(hl[0])
	if errHeader != nil {
		return false, errHeader
	}

	if s.head.MustNot(header) {
		return false, nil
	}

	// 再验证更多字符
	header, errHeader = conn.Pick(hl[1])
	if errHeader != nil {
		return false, errHeader
	}

	if s.head.Is(header) {
		return true, nil
	}
	return false, nil
}

func (s *server) Listener() *listener {
	return s.listener
}

func (s *server) Head() Protocol {
	return s.head
}

func newListener(size int) *listener {
	return &listener{
		Connects: make(chan net.Conn, size),
	}
}

type listener struct {
	AddrValue net.Addr
	Connects  chan net.Conn
}

func (l *listener) SetAddr(addr net.Addr) {
	l.AddrValue = addr
}

func (l *listener) Accept() (net.Conn, error) {
	return <-l.Connects, nil
}

func (l *listener) Close() error {
	close(l.Connects)
	return nil
}

func (l *listener) Addr() net.Addr {
	return l.AddrValue
}

func (l *listener) DispatchConnAsync(conn net.Conn) {
	l.Connects <- conn
}

func newConn(raw net.Conn, opts *Options) *conn {
	return &conn{
		Conn: raw,
		opts: opts,
	}
}

type conn struct {
	net.Conn
	header       []byte
	headerReader io.Reader

	headerHas bool
	opts      *Options
}

// Pick 获取请求头，用于判断协议
func (c *conn) Pick(expectLen int) ([]byte, error) {
	nl := expectLen - len(c.header)
	if nl <= 0 {
		return c.header[:expectLen], nil
	}

	buf := make([]byte, nl)
	n, err := c.Conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != nl {
		return nil, fmt.Errorf("expect read length=%d,got=%d", nl, n)
	}

	c.header = append(c.header, buf...)
	c.headerReader = bytes.NewReader(c.header)
	c.headerHas = true
	return c.header, nil
}

// Read 读取内容
func (c *conn) Read(b []byte) (int, error) {
	return c.read(b)
}

func (c *conn) read(b []byte) (int, error) {
	if !c.headerHas {
		return c.Conn.Read(b)
	}
	n, err := c.headerReader.Read(b)

	if err != nil {
		if err == io.EOF {
			c.headerHas = false
		} else {
			return n, err
		}
	}

	if n == len(b) {
		return n, nil
	}
	m, errM := c.Conn.Read(b[n:])
	return m + n, errM
}

func (c *conn) Close() error {
	c.opts.invokeOnConnClose(c.Conn)
	return c.Conn.Close()
}
