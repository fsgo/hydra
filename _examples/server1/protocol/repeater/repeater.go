/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/29
 */

package repeater

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync/atomic"

	"github.com/fsgo/mpserver/protocol"
)

type Protocol struct {
	config  *protocol.Config
	running int32
	Handler func(input []byte, writer io.Writer)
}

func (p *Protocol) HeaderLen() int {
	return 4
}

func (p *Protocol) Is(header []byte) bool {
	return bytes.HasPrefix(header, []byte("say:"))
}

func (p *Protocol) BindConfig(config *protocol.Config) {
	p.config = config
}

func (p *Protocol) Name() string {
	return "repeater"
}

func (p *Protocol) Serve(l net.Listener) error {
	p.running = 1

	for p.Running() {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go p.handle(conn)
	}
	return nil
}

func (p *Protocol) Running() bool {
	return atomic.LoadInt32(&p.running) == 1
}

func (p *Protocol) handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for p.Running() {
		line, _, err := reader.ReadLine()
		if err != nil {
			conn.Write([]byte("ReadLine with error:"))
			conn.Write([]byte(err.Error()))
			return
		}

		if !bytes.HasPrefix(line, []byte("say:")) {
			conn.Write([]byte("wrong input,expect start with: \"say:\""))
			return
		}
		body := line[4:]
		conn.Write([]byte("reply:"))
		conn.Write(body)
		conn.Write([]byte("\n"))
	}
}

func (p *Protocol) Close() error {
	atomic.StoreInt32(&p.running, 0)
	return nil
}

var _ protocol.Protocol = &Protocol{}
