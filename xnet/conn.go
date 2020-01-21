/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/18
 */

package xnet

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func NewConn(conn net.Conn, opts *Options) Conn {
	return &ConnProxy{
		Conn: conn,
		Opts: opts,
	}
}

type Conn interface {
	net.Conn

	Header(expectLen int) ([]byte, error)
	OnConnect() error
}

type ConnProxy struct {
	net.Conn
	header       []byte
	headerReader io.Reader

	headerHas bool
	Opts      *Options
}

// Header 获取请求头，用于判断协议
func (c *ConnProxy) Header(expectLen int) ([]byte, error) {
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
func (c *ConnProxy) Read(b []byte) (int, error) {
	return c.read(b)
}

func (c *ConnProxy) read(b []byte) (int, error) {
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

func (c *ConnProxy) OnConnect() error {
	if c.Opts.OnConnect != nil {
		return c.Opts.OnConnect(c.Conn)
	}
	return nil
}

func (c *ConnProxy) Close() error {
	if c.Opts.OnConnClose != nil {
		c.Opts.OnConnClose(c.Conn)
	}
	return c.Conn.Close()
}

var _ Conn = &ConnProxy{}
