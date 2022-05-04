// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package hydra

import (
	"log"
	"net"
)

type Options struct {
	ListerChanSize int

	OnConnect func(conn net.Conn) error

	OnConnClose func(conn net.Conn)

	ReadError func(conn net.Conn, err error)

	OnAcceptError func(err error)

	WrongHead func(conn net.Conn)
}

func (o *Options) GetListerChanSize() int {
	if o.ListerChanSize > 0 {
		return o.ListerChanSize
	}
	return 1024
}

func (o *Options) invokeOnConnect(conn net.Conn) error {
	if o.OnConnect == nil {
		return nil
	}
	return o.OnConnect(conn)
}

func (o *Options) invokeOnConnClose(conn net.Conn) {
	if o.OnConnClose == nil {
		return
	}
	o.OnConnClose(conn)
}

func (o *Options) invokeOnAcceptError(err error) {
	if o.OnAcceptError == nil {
		return
	}
	o.OnAcceptError(err)
}

func (o *Options) invokeOnWrongHead(conn net.Conn) {
	if o.WrongHead == nil {
		return
	}
	o.WrongHead(conn)
}

func (o *Options) invokeOnReadError(conn net.Conn, err error) {
	if o.ReadError == nil {
		return
	}
	o.ReadError(conn, err)
}

// DebugOptions 默认用于调试的选型
var DebugOptions = &Options{
	OnConnect: func(conn net.Conn) error {
		log.Println("invokeOnConnect: client=", conn.RemoteAddr().String())
		return nil
	},

	OnConnClose: func(conn net.Conn) {
		log.Println("invokeOnConnClose: client=", conn.RemoteAddr().String())
	},

	OnAcceptError: func(err error) {
		log.Println("invokeOnAcceptError:", err)
	},

	ReadError: func(conn net.Conn, err error) {
		log.Println("invokeOnAcceptError: client=", conn.RemoteAddr(), err)
	},

	WrongHead: func(conn net.Conn) {
		log.Println("WrongHead: client=", conn.RemoteAddr())
	},
}

var optionsEmpty = &Options{}
