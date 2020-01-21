/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/29
 */

package mpserver

import (
	"net"
)

type Listener interface {
	net.Listener
	DispatchConnAsync(conn net.Conn)
}

type ListenerProxy struct {
	addr     net.Addr
	connects chan net.Conn
}

func (l *ListenerProxy) Accept() (net.Conn, error) {
	return <-l.connects, nil
}

func (l *ListenerProxy) Close() error {
	close(l.connects)
	return nil
}

func (l *ListenerProxy) Addr() net.Addr {
	return l.addr
}

func (l *ListenerProxy) DispatchConnAsync(conn net.Conn) {
	l.connects <- conn
}

var _ Listener = &ListenerProxy{}
