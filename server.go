// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/7

package hydra

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// Server 具备优雅关闭的 server
type Server interface {
	// Serve 启动 server
	Serve(net.Listener) error

	// Shutdown 优雅关闭
	Shutdown(context.Context) error
}

var _ Server = (*fixedListenerServer)(nil)

type fixedListenerServer struct {
	Server
	listener net.Listener
}

func (s *fixedListenerServer) Serve(_ net.Listener) error {
	return s.Server.Serve(s.listener)
}

func (s *fixedListenerServer) Shutdown(ctx context.Context) error {
	defer s.listener.Close()
	return s.Server.Shutdown(ctx)
}

// Logger 日志接口定义
type Logger interface {
	Println(v ...any)
	Printf(format string, v ...any)
}

// Starter 优雅关闭 server
type Starter struct {
	Server   Server
	Listener net.Listener

	Timeout time.Duration
	Signals []os.Signal
	Logger  Logger
}

func (sd *Starter) RunGrace() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sd.getSignals()...)
	g := &errgroup.Group{}
	g.Go(func() error {
		return sd.Server.Serve(sd.Listener)
	})
	g.Go(func() error {
		return sd.shutdown(sd.Server, ch)
	})
	return g.Wait()
}

func (sd *Starter) shutdown(ser Server, ch <-chan os.Signal) error {
	sig := <-ch
	sd.getLogger().Println("[Starter] receive signal ", sig, ", start Shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), sd.getTimeout())
	defer cancel()
	start := time.Now()
	err := ser.Shutdown(ctx)
	sd.getLogger().Println("[Starter] Shutdown finish, err=", err, ", cost=", time.Since(start))
	return err

}

func (sd *Starter) getSignals() []os.Signal {
	if len(sd.Signals) == 0 {
		return []os.Signal{os.Interrupt, syscall.SIGTERM}
	}
	return sd.Signals
}

func (sd *Starter) getTimeout() time.Duration {
	if sd.Timeout > 0 {
		return sd.Timeout
	}
	return time.Minute
}

func (sd *Starter) getLogger() Logger {
	if sd.Logger != nil {
		return sd.Logger
	}
	return log.Default()
}

func RunGrace(ser Server, l net.Listener, timeout time.Duration) error {
	gs := &Starter{
		Timeout:  timeout,
		Listener: l,
		Server:   ser,
	}
	return gs.RunGrace()
}
