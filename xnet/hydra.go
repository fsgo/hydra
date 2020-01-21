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

	"github.com/fsgo/hydra/servers"
)

var ErrUnknownProtocol = fmt.Errorf("protocol not support yet")

type Hydra interface {
	RegisterServer(p servers.Server)
	Listen(addr net.Addr) (l net.Listener, err error)
	Serve(listener net.Listener, errChan chan error)
	Dispatch(conn net.Conn) error
	Stop() error
}

func NewHydra(opts *Options) Hydra {

	if opts == nil {
		opts = OptionsDefault
	}

	return &hydra{
		Opts: opts,
	}
}

type hydra struct {
	Opts    *Options
	Servers []Server
}

func (ps *hydra) RegisterServer(p servers.Server) {
	ps.Servers = append(ps.Servers, newServer(p, ps.Opts.GetListerChanSize()))

	sort.Slice(ps.Servers, func(i, j int) bool {
		a := ps.Servers[i]
		b := ps.Servers[j]
		return a.Server().HeaderLen() > b.Server().HeaderLen()
	})
}

func (ps *hydra) Listen(addr net.Addr) (l net.Listener, err error) {
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

func (ps *hydra) Serve(listener net.Listener, errChan chan error) {
	for _, server := range ps.Servers {
		server.Listener().SetAddr(listener.Addr())
		go func(ls Server) {
			errChan <- ls.Serve()
		}(server)
	}
}

func (ps *hydra) Dispatch(conn net.Conn) error {
	myConn := NewConn(conn, ps.Opts)
	if err := myConn.OnConnect(); err != nil {
		return err
	}

	var header []byte
	var errHeader error

	for _, server := range ps.Servers {
		header, errHeader = myConn.Header(server.Server().HeaderLen())
		if errHeader != nil {
			conn.Write([]byte(errHeader.Error()))
			myConn.Close()
			return errHeader
		}

		if server.Server().Is(header) {
			server.Listener().DispatchConnAsync(myConn)
			return nil
		}
	}

	myConn.Close()
	return ErrUnknownProtocol
}

func (ps *hydra) Stop() error {
	for _, p := range ps.Servers {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}

var _ Hydra = &hydra{}
