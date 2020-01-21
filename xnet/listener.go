/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/18
 */

package xnet

import (
	"net"
)

type Listener interface {
	net.Listener
	SetAddr(addr net.Addr)
	DispatchConnAsync(conn net.Conn)
}

func NewListener(size int) Listener {
	return &ListenerProxy{
		Connects: make(chan net.Conn, size),
	}
}

type ListenerProxy struct {
	AddrValue net.Addr
	Connects  chan net.Conn
}

func (l *ListenerProxy) SetAddr(addr net.Addr) {
	l.AddrValue = addr
}

func (l *ListenerProxy) Accept() (net.Conn, error) {
	return <-l.Connects, nil
}

func (l *ListenerProxy) Close() error {
	close(l.Connects)
	return nil
}

func (l *ListenerProxy) Addr() net.Addr {
	return l.AddrValue
}

func (l *ListenerProxy) DispatchConnAsync(conn net.Conn) {
	l.Connects <- conn
}

var _ Listener = &ListenerProxy{}
