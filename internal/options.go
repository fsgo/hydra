// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package internal

import (
	"log"
	"net"
)

type Options struct {
	ListerChanSize int

	FnOnConnect func(conn net.Conn) error

	FnOnConnClose func(conn net.Conn)

	FnReadError func(conn net.Conn, err error)

	FnOnAcceptError func(err error)

	FnWrongHead func(conn net.Conn)
}

func (o *Options) GetListerChanSize() int {
	if o.ListerChanSize > 0 {
		return o.ListerChanSize
	}
	return 1024
}

func (o *Options) OnConnect(conn net.Conn) error {
	if o.FnOnConnect == nil {
		return nil
	}
	return o.FnOnConnect(conn)
}

func (o *Options) OnConnClose(conn net.Conn) {
	if o.FnOnConnClose == nil {
		return
	}
	o.FnOnConnClose(conn)
}

func (o *Options) OnAcceptError(err error) {
	if o.FnOnAcceptError == nil {
		return
	}
	o.FnOnAcceptError(err)
}

func (o *Options) OnWrongHead(conn net.Conn) {
	if o.FnWrongHead == nil {
		return
	}
	o.FnWrongHead(conn)
}

func (o *Options) OnReadError(conn net.Conn, err error) {
	if o.FnReadError == nil {
		return
	}
	o.FnReadError(conn, err)
}

// OptionsDebug 默认用于调试的选型
var OptionsDebug = &Options{

	FnOnConnect: func(conn net.Conn) error {
		log.Println("OnConnect: client=", conn.RemoteAddr().String())
		return nil
	},

	FnOnConnClose: func(conn net.Conn) {
		log.Println("OnConnClose: client=", conn.RemoteAddr().String())
	},

	FnOnAcceptError: func(err error) {
		log.Println("OnAcceptError:", err)
	},

	FnReadError: func(conn net.Conn, err error) {
		log.Println("OnAcceptError: client=", conn.RemoteAddr(), err)
	},

	FnWrongHead: func(conn net.Conn) {
		log.Println("WrongHead: client=", conn.RemoteAddr())
	},
}

var OptionsEmpty = &Options{}
