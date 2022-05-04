// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package hydra

import (
	"context"
	"net"
	"sort"
	"sync/atomic"
)

type Hydra struct {
	Opts    *Options
	servers []*server
	running int32
}

func (h *Hydra) getOpts() *Options {
	if h.Opts == nil {
		return optionsEmpty
	}
	return h.Opts
}

// Serve 监听服务
// 注意，必须在所有 BindHead 完成之后
func (h *Hydra) Serve(listener net.Listener) error {
	atomic.StoreInt32(&h.running, 1)
	for {
		if atomic.LoadInt32(&h.running) != 1 {
			break
		}
		c, err := listener.Accept()

		if err != nil {
			h.getOpts().invokeOnAcceptError(err)
			continue
		}

		go h.dispatch(c)
	}
	return nil
}

// BindHead 绑定一个协议头，并返回对应的 Listener
func (h *Hydra) BindHead(p Protocol) (ls net.Listener, err error) {
	ln := newListener(h.getOpts().ListerChanSize)
	ser := newServer(p, ln)
	h.servers = append(h.servers, ser)

	sort.Slice(h.servers, func(i, j int) bool {
		a := h.servers[i]
		b := h.servers[j]
		return a.Head().DiscernLengths()[1] > b.Head().DiscernLengths()[1]
	})
	return ln, err
}

func (h *Hydra) dispatch(conn net.Conn) {
	xc := newConn(conn, h.getOpts())

	if err := h.getOpts().invokeOnConnect(xc); err != nil {
		xc.Close()
		return
	}

	for _, p := range h.servers {
		is, err := p.Is(xc)
		if err != nil {
			h.getOpts().invokeOnReadError(conn, err)
			xc.Close()
			return
		}
		if is {
			p.Listener().DispatchConnAsync(xc)
			return
		}
	}

	// 不识别的协议
	h.getOpts().invokeOnWrongHead(xc)
	xc.Close()
}

// Shutdown 停止服务
func (h *Hydra) Shutdown(ctx context.Context) error {
	atomic.StoreInt32(&h.running, 0)
	return nil
}
