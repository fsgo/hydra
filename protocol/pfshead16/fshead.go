/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/12
 */

package pfshead16

import (
	"net"
	"sync/atomic"

	"github.com/fsgo/fsprotocol/fshead16"

	"github.com/fsgo/hydra/protocol"
)

var headLen protocol.DiscernLengths

func init() {
	headLen = [2]int{
		fshead16.DiscernLength,
		fshead16.DiscernLength,
	}
}

// protocol HTTP 协议
type FSHead16 struct {
	Head    fshead16.Head
	Handler func(conn net.Conn)
	running int32
}

func (p *FSHead16) MustNot(header []byte) bool {
	return false
}

func (p *FSHead16) HeaderLen() protocol.DiscernLengths {
	return headLen
}

func (p *FSHead16) Is(header []byte) bool {
	return p.Head.Is(header)
}

func (p *FSHead16) Name() string {
	return "fshead16"
}

func (p *FSHead16) Serve(l net.Listener) error {
	p.running = 1
	for {
		if atomic.LoadInt32(&p.running) != 1 {
			return nil
		}

		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go p.Handler(conn)
	}
}

func (p *FSHead16) Close() error {
	atomic.StoreInt32(&p.running, 0)
	return nil
}

var _ protocol.Protocol = &FSHead16{}
