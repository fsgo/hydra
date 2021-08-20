// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package hydra

import (
	"net"
	"sort"
	"sync/atomic"

	"github.com/fsgo/hydra/internal"
)

// Head 协议识别的接口定义
type Head = internal.Head

type DiscernLengths = internal.DiscernLengths

type Options = internal.Options

var OptionsEmpty = internal.OptionsEmpty

var OptionsDebug = internal.OptionsDebug

// Hydra 识别多协议
type Hydra interface {
	// BindHead 绑定一个协议头，并返回对应的Listener
	BindHead(p Head) (ls net.Listener, err error)

	// Serve 监听服务
	// 注意，必须在所有 BindHead 完成之后
	Serve(listener net.Listener) error

	// Stop 停止服务
	Stop() error
}

func New(opts *Options) Hydra {
	if opts == nil {
		opts = internal.OptionsEmpty
	}
	return &hydra{
		servers: nil,
		opts:    opts,
	}
}

type hydra struct {
	opts    *Options
	servers []internal.XServer
	running int32
}

func (h *hydra) Serve(listener net.Listener) error {
	atomic.StoreInt32(&h.running, 1)
	for {
		if atomic.LoadInt32(&h.running) != 1 {
			break
		}
		conn, err := listener.Accept()

		if err != nil {
			h.opts.OnAcceptError(err)
			continue
		}

		go h.dispatch(conn)
	}
	return nil
}

func (h *hydra) BindHead(p Head) (ls net.Listener, err error) {
	ln := internal.NewListener(h.opts.ListerChanSize)
	ser := internal.NewServer(p, ln)
	h.servers = append(h.servers, ser)

	sort.Slice(h.servers, func(i, j int) bool {
		a := h.servers[i]
		b := h.servers[j]
		return a.Head().HeaderLen()[1] > b.Head().HeaderLen()[1]
	})
	return ln, err
}

func (h *hydra) dispatch(conn net.Conn) {
	xc := internal.NewConn(conn, h.opts)

	if err := h.opts.OnConnect(xc); err != nil {
		xc.Close()
		return
	}

	for _, p := range h.servers {
		is, err := p.Is(xc)
		if err != nil {
			h.opts.OnReadError(conn, err)
			xc.Close()
			return
		}
		if is {
			p.Listener().DispatchConnAsync(xc)
			return
		}
	}

	// 不识别的协议
	h.opts.OnWrongHead(xc)
	xc.Close()
}

func (h *hydra) Stop() error {
	atomic.StoreInt32(&h.running, 0)
	return nil
}

var _ Hydra = &hydra{}
