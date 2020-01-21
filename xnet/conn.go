/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/29
 */

package mpserver

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func newConn(cn net.Conn, opts *Options) *conn {
	return &conn{
		Conn: cn,
		Opts: opts,
	}
}

type conn struct {
	net.Conn
	header       []byte
	headerReader io.Reader

	headerHas bool
	Opts      *Options
}

// Header 获取请求头，用于判断协议
func (c *conn) Header(expectLen int) ([]byte, error) {
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

func (c *conn) OnConnect() error {
	if c.Opts.OnConnect != nil {
		return c.Opts.OnConnect(c.Conn)
	}
	return nil
}

func (c *conn) Close() error {
	if c.Opts.OnConnClose != nil {
		c.Opts.OnConnClose(c.Conn)
	}
	return c.Conn.Close()
}
