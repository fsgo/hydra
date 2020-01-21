/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/18
 */

package xnet

import (
	"net"

	"github.com/fsgo/hydra/protocol"
)

type Listener interface {
	net.Listener
	SetAddr(addr net.Addr)
	Protocol() protocol.Protocol
	DispatchConnAsync(conn net.Conn)
	Serve() error
}

func NewListener(p protocol.Protocol) Listener {
	return &ListenerProxy{
		ProtocolValue: p,
		Connects:      make(chan net.Conn, 1024),
	}
}

type ListenerProxy struct {
	AddrValue     net.Addr
	Connects      chan net.Conn
	ProtocolValue protocol.Protocol
}

func (l *ListenerProxy) Protocol() protocol.Protocol {
	return l.ProtocolValue
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

func (l *ListenerProxy) Serve() error {
	return l.ProtocolValue.Serve(l)
}

var _ Listener = &ListenerProxy{}
