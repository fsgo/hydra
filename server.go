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
)

// Server 具备优雅关闭的 server
type Server interface {
	// Serve 启动 server
	Serve(net.Listener) error

	// Shutdown 优雅关闭
	Shutdown(context.Context) error
}

type CanShutdown interface {
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
	return s.Server.Shutdown(ctx)
}

// Logger 日志接口定义
type Logger interface {
	Println(v ...any)
	Printf(format string, v ...any)
}

// GraceShutdown 优雅关闭 server
type GraceShutdown struct {
	Timeout time.Duration
	Signals []os.Signal
	Logger  Logger
}

// WaitSignal 异步等待退出信号
func (sd *GraceShutdown) WaitSignal(ser CanShutdown) {
	go sd.wait(ser)
}

func (sd *GraceShutdown) wait(ser CanShutdown) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sd.getSignals()...)
	sig := <-ch

	sd.getLogger().Println("[GraceShutdown] receive signal ", sig, ", start Shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), sd.getTimeout())
	defer cancel()
	err := ser.Shutdown(ctx)
	sd.getLogger().Println("[GraceShutdown] Shutdown finish, err=", err)
}

func (sd *GraceShutdown) getSignals() []os.Signal {
	if len(sd.Signals) == 0 {
		return []os.Signal{os.Interrupt, syscall.SIGTERM}
	}
	return sd.Signals
}

func (sd *GraceShutdown) getTimeout() time.Duration {
	if sd.Timeout > 0 {
		return sd.Timeout
	}
	return time.Minute
}

func (sd *GraceShutdown) getLogger() Logger {
	if sd.Logger != nil {
		return sd.Logger
	}
	return log.Default()
}

func WaitShutdown(ser CanShutdown, timeout time.Duration) {
	gs := &GraceShutdown{
		Timeout: timeout,
	}
	gs.WaitSignal(ser)
}
