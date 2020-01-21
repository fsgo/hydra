/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/12
 */

package fshead16server

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/fsgo/fsprotocol/fshead16"

	"github.com/fsgo/hydra/protocol"
)

type Read func() ([]byte, error)

// Protocol HTTP 协议
type Protocol struct {
	Head    fshead16.Head
	Handler func(conn net.Conn)
	running int32
}

func (p *Protocol) HeaderLen() int {
	return p.Head.DiscernLen()
}

func (p *Protocol) Is(header []byte) bool {
	return p.Head.Is(header)
}

func (p *Protocol) Name() string {
	return "fshead16"
}

func (p *Protocol) Serve(l net.Listener) error {
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

var headerBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, fshead16.Length)
	},
}

func (p *Protocol) Close() error {
	atomic.StoreInt32(&p.running, 0)
	return nil
}

var _ protocol.Protocol = &Protocol{}
