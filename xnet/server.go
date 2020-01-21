/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/18
 */

package xnet

import (
	"errors"
	"fmt"
	"net"
	"sort"

	"github.com/fsgo/mpserver/protocol"
)

var ErrUnknownProtocol = fmt.Errorf("protocol not support yet")

type Server interface {
	RegisterProtocol(p protocol.Protocol)
	Listen(addr net.Addr) (l net.Listener, err error)
	Serve(listener net.Listener, errChan chan error)
	Dispatch(conn net.Conn) error
	Stop() error
}

func NewServer(opts *Options) Server {
	return &XServer{
		Opts: opts,
	}
}

type XServer struct {
	Opts      *Options
	Listeners []Listener
}

func (ps *XServer) RegisterProtocol(p protocol.Protocol) {
	ps.Listeners = append(ps.Listeners, NewListener(p))

	sort.Slice(ps.Listeners, func(i, j int) bool {
		a := ps.Listeners[i]
		b := ps.Listeners[j]
		return a.Protocol().HeaderLen() > b.Protocol().HeaderLen()
	})
}

func (ps *XServer) Listen(addr net.Addr) (l net.Listener, err error) {
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

func (ps *XServer) Serve(listener net.Listener, errChan chan error) {
	for _, ls := range ps.Listeners {
		ls.SetAddr(listener.Addr())
		go func(ls Listener) {
			errChan <- ls.Serve()
		}(ls)
	}
}

func (ps *XServer) Dispatch(conn net.Conn) error {
	myConn := NewConn(conn, ps.Opts)
	if err := myConn.OnConnect(); err != nil {
		return err
	}

	var header []byte
	var errHeader error

	for _, ls := range ps.Listeners {
		header, errHeader = myConn.Header(ls.Protocol().HeaderLen())
		if errHeader != nil {
			conn.Write([]byte(errHeader.Error()))
			myConn.Close()
			return errHeader
		}

		if ls.Protocol().Is(header) {
			ls.DispatchConnAsync(myConn)
			return nil
		}
	}

	myConn.Close()
	return ErrUnknownProtocol
}

func (ps *XServer) Stop() error {
	for _, p := range ps.Listeners {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}

var _ Server = &XServer{}
