/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package hydra

import (
	"errors"
	"net"
	"sync/atomic"

	"github.com/fsgo/hydra/protocol"
	"github.com/fsgo/hydra/xnet"
)

type Options = xnet.Options

// HydraServer 多协议server接口定义
type Server interface {
	SetListenAddr(addr net.Addr)

	RegisterProtocol(s protocol.Protocol)

	Start() error

	Stop() error
}

// NewServer 一个新的server
func NewServer(opts *Options) Server {
	if opts == nil {
		opts = xnet.OptionsEmpty
	}
	return &defaultServer{
		opts:  opts,
		hydra: xnet.NewHydra(opts),
	}
}

// defaultServer 多协议server的默认实现
type defaultServer struct {
	addr net.Addr
	opts *xnet.Options

	hydra xnet.Hydra

	running int32
}

// SetListenAddr 设置监听的地址
func (s *defaultServer) SetListenAddr(addr net.Addr) {
	s.addr = addr
}

// RegisterProtocol 注册一种新协议
func (s *defaultServer) RegisterProtocol(ss protocol.Protocol) {
	s.hydra.RegisterProtocol(ss)
}

// Start 启动服务
func (s *defaultServer) Start() error {
	s.running = 1
	if s.addr == nil {
		return errors.New("addr is nil")
	}

	listener, err := s.hydra.Listen(s.addr)

	if err != nil {
		return err
	}

	defer listener.Close()

	errChan := make(chan error, 1)

	s.hydra.Serve(listener, errChan)

	for {
		if atomic.LoadInt32(&s.running) != 1 {
			break
		}

		select {
		case err := <-errChan:
			return err
		default:
		}

		conn, err := listener.Accept()

		if err != nil {
			s.opts.OnAcceptError(err)
			continue
		}

		go s.hydra.Dispatch(conn)
	}

	return errors.New("server exit")
}

// Stop 停止服务
func (s *defaultServer) Stop() error {
	if err := s.hydra.Stop(); err != nil {
		return err
	}
	atomic.StoreInt32(&s.running, 0)
	return nil
}

var _ Server = &defaultServer{}
