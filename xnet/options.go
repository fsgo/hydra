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
	OnListen       func(l net.Listener) error
	OnConnect      func(conn net.Conn) error
	OnConnClose    func(conn net.Conn)

	OnAcceptError func(err error)
}

func (o *Options) GetListerChanSize() int {
	if o.ListerChanSize > 0 {
		return o.ListerChanSize
	}
	return 1024
}

var OptionsDefault = &Options{
	OnListen: func(l net.Listener) error {
		log.Println("OnListen:", l.Addr())
		return nil
	},

	OnConnect: func(conn net.Conn) error {
		log.Println("OnConnect,client=", conn.RemoteAddr().String())
		return nil
	},

	OnConnClose: func(conn net.Conn) {
		log.Println("OnConnClose", conn.RemoteAddr().String())
	},
}

var OptionsNil = &Options{}
