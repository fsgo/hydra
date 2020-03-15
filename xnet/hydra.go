/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/18
 */

package xnet

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sort"

	"github.com/fsgo/hydra/protocol"
)

var ErrUnknownProtocol = fmt.Errorf("protocol not support yet")

type Hydra interface {
	RegisterProtocol(p protocol.Protocol)
	Listen(addr net.Addr) (l net.Listener, err error)
	Serve(listener net.Listener, errChan chan error)
	Dispatch(conn net.Conn) error
	Stop() error
}

func NewHydra(opts *Options) Hydra {
	if opts == nil {
		opts = OptionsEmpty
	}

	return &hydra{
		Opts: opts,
	}
}

type hydra struct {
	Opts    *Options
	Servers []Server
}

func (ps *hydra) RegisterProtocol(p protocol.Protocol) {
	ps.Servers = append(ps.Servers, newServer(p, ps.Opts.GetListerChanSize()))

	sort.Slice(ps.Servers, func(i, j int) bool {
		a := ps.Servers[i]
		b := ps.Servers[j]
		return a.Protocol().HeaderLen()[1] > b.Protocol().HeaderLen()[1]
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
	if err != nil {
		return l, err
	}

	if errCallBack := ps.Opts.OnListen(l); errCallBack != nil {
		l.Close()
		return nil, errCallBack
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
	if err := ps.Opts.OnConnect(myConn); err != nil {
		return err
	}

	for _, server := range ps.Servers {
		is, err := server.Is(myConn)
		if err != nil {
			conn.Write([]byte(err.Error()))
			myConn.Close()
			return err
		}

		if is {
			server.Listener().DispatchConnAsync(myConn)
			return nil
		}
	}

	myConn.Close()
	return ErrUnknownProtocol
}

func (ps *hydra) Stop() error {
	var buf bytes.Buffer
	for _, p := range ps.Servers {
		if err := p.Close(); err != nil {
			buf.WriteString(fmt.Sprintf("server=%s Close with error:=%v;", p.Protocol().Name(), err))
		}
	}
	if buf.Len() == 0 {
		return nil
	}
	return errors.New(buf.String())
}

var _ Hydra = &hydra{}
