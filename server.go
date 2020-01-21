/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package hydra

import (
	"errors"
	"log"
	"net"
	"sync/atomic"

	"github.com/fsgo/hydra/protocol"
	"github.com/fsgo/hydra/xnet"
)

type Options = xnet.Options

// MultiProtocolServer 多协议server接口定义
type MultiProtocolServer interface {
	SetListenAddr(addr net.Addr)

	RegisterProtocol(p protocol.Protocol)

	Start() error

	Stop() error
}

// NewServer 一个新的server
func NewServer(opts *Options) MultiProtocolServer {
	if opts == nil {
		opts = xnet.OptionsDefault
	}
	return &Server{
		xServer: xnet.NewServer(opts),
	}
}

// Server 多协议server的默认实现
type Server struct {
	addr net.Addr

	xServer xnet.Server

	running int32
}

// SetListenAddr 设置监听的地址
func (s *Server) SetListenAddr(addr net.Addr) {
	s.addr = addr
}

// RegisterProtocol 注册一种新协议
func (s *Server) RegisterProtocol(p protocol.Protocol) {
	s.xServer.RegisterProtocol(p)
}

// Start 启动服务
func (s *Server) Start() error {
	s.running = 1
	if s.addr == nil {
		return errors.New("addr is nil")
	}

	listener, err := s.xServer.Listen(s.addr)

	if err != nil {
		return err
	}

	defer listener.Close()

	log.Printf("server Lister : %s/%s\n", s.addr.Network(), s.addr.String())

	errChan := make(chan error, 1)

	s.xServer.Serve(listener, errChan)

	for {
		if atomic.LoadInt32(&s.running) != 1 {
			break
		}

		select {
		case err := <-errChan:
			log.Println("Serve error:", err)
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("listener.Accept error:", err)
		}
		go s.xServer.Dispatch(conn)
	}

	return errors.New("server exists")
}

// Stop 停止服务
func (s *Server) Stop() error {
	if err := s.xServer.Stop(); err != nil {
		return err
	}
	atomic.StoreInt32(&s.running, 0)
	return nil
}

var _ MultiProtocolServer = &Server{}
