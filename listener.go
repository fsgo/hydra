/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/29
 */

package mpserver

import (
	"net"
)

type listenerProxy struct {
	addr  net.Addr
	conns chan net.Conn
}

func (l *listenerProxy) Accept() (net.Conn, error) {
	return <-l.conns, nil
}

func (l *listenerProxy) Close() error {
	close(l.conns)
	return nil
}

func (l *listenerProxy) Addr() net.Addr {
	return l.addr
}

func (l *listenerProxy) DispatchConnAsync(conn net.Conn) {
	l.conns <- conn
}

var _ net.Listener = &listenerProxy{}
