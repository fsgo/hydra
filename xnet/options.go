/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/18
 */

package xnet

import (
	"log"
	"net"
)

type Options struct {
	ListerChanSize int
	OnListenFn     func(l net.Listener) error
	OnConnectFn    func(conn net.Conn) error
	OnConnCloseFn  func(conn net.Conn)

	OnAcceptErrorFn func(err error)
}

func (o *Options) GetListerChanSize() int {
	if o.ListerChanSize > 0 {
		return o.ListerChanSize
	}
	return 1024
}

func (o *Options) OnListen(l net.Listener) error {
	if o.OnListenFn == nil {
		return nil
	}
	return o.OnListenFn(l)
}

func (o *Options) OnConnect(conn net.Conn) error {
	if o.OnConnectFn == nil {
		return nil
	}
	return o.OnConnect(conn)
}
func (o *Options) OnConnClose(conn net.Conn) {
	if o.OnConnCloseFn == nil {
		return
	}
	o.OnConnCloseFn(conn)
}
func (o *Options) OnAcceptError(err error) {
	if o.OnAcceptErrorFn == nil {
		return
	}
	o.OnAcceptErrorFn(err)
}

var OptionsDefault = &Options{
	OnListenFn: func(l net.Listener) error {
		log.Println("OnListen:", l.Addr())
		return nil
	},

	OnConnectFn: func(conn net.Conn) error {
		log.Println("OnConnect,client=", conn.RemoteAddr().String())
		return nil
	},

	OnConnCloseFn: func(conn net.Conn) {
		log.Println("OnConnClose", conn.RemoteAddr().String())
	},

	OnAcceptErrorFn: func(err error) {
		log.Println("OnAcceptError", err)
	},
}

var OptionsEmpty = &Options{}
