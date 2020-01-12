/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/29
 */

package mpserver

import (
	"errors"
	"fmt"
	"net"
	"sort"

	"github.com/fsgo/mpserver/protocol"
)

func newProtocol(p protocol.Protocol) *serverProtocol {
	return &serverProtocol{
		Protocol: p,
		ListenerProxy: &listenerProxy{
			conns: make(chan net.Conn, 100),
		},
	}
}

type serverProtocol struct {
	Protocol      protocol.Protocol
	ListenerProxy *listenerProxy
}

func (sp *serverProtocol) Close() error {
	if err := sp.ListenerProxy.Close(); err != nil {
		return err
	}
	return nil
}

type protocols struct {
	Opts      *Options
	protocols []*serverProtocol
}

func (ps *protocols) Register(p protocol.Protocol) {
	ps.protocols = append(ps.protocols, newProtocol(p))
	sort.Slice(ps.protocols, func(i, j int) bool {
		a := ps.protocols[i]
		b := ps.protocols[j]
		return a.Protocol.HeaderLen() > b.Protocol.HeaderLen()
	})
}

func (ps *protocols) Listen(addr net.Addr) (l net.Listener, err error) {
	switch addr.(type) {
	case *net.TCPAddr:
		l, err = net.ListenTCP(addr.Network(), addr.(*net.TCPAddr))
	case *net.UnixAddr:
		l, err = net.ListenUnix(addr.Network(), addr.(*net.UnixAddr))
	default:
		l, err = nil, errors.New("not support addr:"+addr.String())
	}

	if l != nil && ps.Opts.OnListen != nil {
		if errCallBack := ps.Opts.OnListen(l); errCallBack != nil {
			l.Close()
			return nil, errCallBack
		}
	}

	return l, err
}

func (ps *protocols) Serve(listener net.Listener, errChan chan error) {
	for _, p := range ps.protocols {
		p.ListenerProxy.addr = listener.Addr()
		go func(p *serverProtocol) {
			errChan <- p.Protocol.Serve(p.ListenerProxy)
		}(p)
	}
}

func (ps *protocols) Dispatch(conn net.Conn) error {
	myConn := newConn(conn, ps.Opts)
	if err := myConn.OnConnect(); err != nil {
		return err
	}

	var header []byte
	var errHeader error

	for _, p := range ps.protocols {
		header, errHeader = myConn.Header(p.Protocol.HeaderLen())
		if errHeader != nil {
			conn.Write([]byte(errHeader.Error()))
			myConn.Close()
			return errHeader
		}

		if p.Protocol.Is(header) {
			p.ListenerProxy.DispatchConnAsync(myConn)
			return nil
		}
	}
	return fmt.Errorf("protocol not support yet")
}

func (ps *protocols) Stop() error {
	for _, p := range ps.protocols {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}
