/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/3/29
 */

package internal

import (
	"net"
)

type XListener interface {
	net.Listener
	SetAddr(addr net.Addr)
	DispatchConnAsync(conn net.Conn)
}

func NewListener(size int) XListener {
	return &Listener{
		Connects: make(chan net.Conn, size),
	}
}

type Listener struct {
	AddrValue net.Addr
	Connects  chan net.Conn
}

func (l *Listener) SetAddr(addr net.Addr) {
	l.AddrValue = addr
}

func (l *Listener) Accept() (net.Conn, error) {
	return <-l.Connects, nil
}

func (l *Listener) Close() error {
	close(l.Connects)
	return nil
}

func (l *Listener) Addr() net.Addr {
	return l.AddrValue
}

func (l *Listener) DispatchConnAsync(conn net.Conn) {
	l.Connects <- conn
}

var _ XListener = &Listener{}
