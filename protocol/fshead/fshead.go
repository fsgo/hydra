/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/12
 */

package fshead

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fsgo/fshead"

	"github.com/fsgo/mpserver/protocol"
)

type Read func() ([]byte, error)

// Protocol HTTP 协议
type Protocol struct {
	Config   *protocol.Config
	MagicNum uint32

	Handler func(h *fshead.FsHead, metaRead Read, bodyRead Read) []byte
}

func (p *Protocol) HeaderLen() int {
	return fshead.Length
}

func (p *Protocol) Is(header []byte) bool {
	_, err := fshead.ParserBytes(header, p.MagicNum)
	return err == nil
}

func (p *Protocol) BindConfig(config *protocol.Config) {
	p.Config = config
}

func (p *Protocol) Name() string {
	return "fshead"
}

func (p *Protocol) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go p.serveConn(conn)
	}
}

var headerBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, fshead.Length)
	},
}

func (p *Protocol) serveConn(conn net.Conn) {
	defer conn.Close()

	readFn := func(buf []byte, wantLen int) error {
		if p.Config.ReadTimeout != 0 {
			conn.SetReadDeadline(time.Now().Add(p.Config.ReadTimeout))
		}
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		if n != wantLen {
			return fmt.Errorf("read %d bytes, want %d bytes", n, wantLen)
		}
		return nil
	}

	readFn1 := func(n uint32) Read {
		return func() (bytes []byte, err error) {
			if n == 0 {
				return nil, nil
			}
			buf := make([]byte, 0, n)
			if err := readFn(buf, int(n)); err != nil {
				return nil, err
			}
			return buf, nil
		}
	}

	for {
		head := headerBufPool.Get().([]byte)
		if err := readFn(head, fshead.Length); err != nil {
			log.Println(err.Error())
			return
		}
		h, err := fshead.ParserBytes(head, p.MagicNum)

		head = head[:0]
		headerBufPool.Put(head)

		if err != nil {
			log.Println(err)
			return
		}

		resp := p.Handler(h, readFn1(uint32(h.MetaLen)), readFn1(h.BodyLen))

		wrote, errWrote := conn.Write(resp)
		if errWrote != nil {
			log.Printf("errWrote:%s\n", errWrote.Error())
			return
		}
		if wrote != len(resp) {
			log.Printf("wrote %d bytes,want %d bytes\n", wrote, len(resp))
			return
		}
	}
}

func (p *Protocol) Close() error {
	return nil
}

var _ protocol.Protocol = &Protocol{}
