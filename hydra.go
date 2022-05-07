// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package hydra

import (
	"context"
	"net"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

var _ Server = (*Hydra)(nil)

// Hydra 识别多协议的 fixedListenerServer
type Hydra struct {
	// Opts 配置选项，可选
	Opts *Options

	// DefaultListener 未识别的协议，可选
	// 若未设置，在运行时会自动初始化
	DefaultListener Listener

	referees referees
	listener net.Listener
	running  int32

	servers []Server
}

func (h *Hydra) getOpts() *Options {
	if h.Opts == nil {
		return optionsEmpty
	}
	return h.Opts
}

func (h *Hydra) initListeners() {
	if h.DefaultListener == nil {
		h.DefaultListener = NewListener(h.getOpts().ListerChanSize)
		h.DefaultListener.SetLocalAddr(h.listener.Addr())
	}
	for i := 0; i < len(h.referees); i++ {
		r := h.referees[i]
		r.listener.SetLocalAddr(h.listener.Addr())
	}
}

// Listen 绑定一个协议头，并返回对应的 Listener
func (h *Hydra) Listen(p Protocol) net.Listener {
	ln := NewListener(h.getOpts().ListerChanSize)
	h.BindListener(p, ln)
	return ln
}

// BindListener 绑定协议和 Listener
// 解析成功的 Conn 会发送到该 Listener
func (h *Hydra) BindListener(p Protocol, ln Listener) {
	ser := newReferee(p, ln)
	h.referees = append(h.referees, ser)
	h.referees.MustCheck()
	h.referees.Sort()
}

// GetDefaultListener 未知协议的连接
func (h *Hydra) GetDefaultListener() net.Listener {
	return h.DefaultListener
}

// BindServer 绑定协议和 Server
func (h *Hydra) BindServer(p Protocol, ser Server) {
	ns := &fixedListenerServer{
		Server:   ser,
		listener: h.Listen(p),
	}
	h.servers = append(h.servers, ns)
}

// SetDefaultServer 绑定未知协议对应的 Server
func (h *Hydra) SetDefaultServer(ser Server) {
	ns := &fixedListenerServer{
		Server:   ser,
		listener: h.GetDefaultListener(),
	}
	h.servers = append(h.servers, ns)
}

// Serve 监听服务
// 注意，必须在所有 Listen 完成之后
func (h *Hydra) Serve(listener net.Listener) error {
	h.listener = listener
	h.initListeners()
	atomic.StoreInt32(&h.running, 1)

	g := &errgroup.Group{}
	g.Go(func() error {
		return h.serve()
	})
	for i := 0; i < len(h.servers); i++ {
		ser := h.servers[i]
		g.Go(func() error {
			return ser.Serve(nil)
		})
	}
	return g.Wait()
}

func (h *Hydra) serve() error {
	var timeout time.Duration
	var sd canSetDeadline
	var ok bool
	if sd, ok = h.listener.(canSetDeadline); ok {
		timeout = h.getOpts().GetAcceptTimeout()
	}

	for {
		if atomic.LoadInt32(&h.running) != 1 {
			break
		}
		if timeout > 0 {
			sd.SetDeadline(time.Now().Add(timeout))
		}
		c, err := h.listener.Accept()
		if err != nil {
			if os.IsTimeout(err) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			h.getOpts().invokeOnAcceptError(err)
			continue
		}
		go h.dispatch(c)
	}
	return net.ErrClosed
}

func (h *Hydra) dispatch(conn net.Conn) {
	pc := newConn(conn)
	for _, ref := range h.referees {
		is, err := ref.Is(pc)
		if err != nil {
			pc.Close()
			return
		}
		if is {
			ref.listener.Dispatch(pc)
			return
		}
	}
	// 不识别的协议
	h.DefaultListener.Dispatch(pc)
}

// Shutdown 停止服务
func (h *Hydra) Shutdown(ctx context.Context) error {
	atomic.StoreInt32(&h.running, 0)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < len(h.servers); i++ {
		ser := h.servers[i]
		g.Go(func() error {
			return ser.Shutdown(ctx)
		})
	}

	// g.Go(func() error {
	// 	<-ctx.Done()
	// 	h.referees.Close()
	// 	return nil
	// })

	return g.Wait()
}
