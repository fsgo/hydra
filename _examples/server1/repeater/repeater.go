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

	"github.com/fsgo/hydra/servers"
)

type Repeater struct {
	running int32
	Handler func(input []byte, writer io.Writer)
}

func (p *Repeater) HeaderLen() int {
	return 4
}

func (p *Repeater) Is(header []byte) bool {
	return bytes.HasPrefix(header, []byte("say:"))
}

func (p *Repeater) Name() string {
	return "repeater"
}

func (p *Repeater) Serve(l net.Listener) error {
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

func (p *Repeater) Running() bool {
	return atomic.LoadInt32(&p.running) == 1
}

func (p *Repeater) handle(conn net.Conn) {
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
			conn.Write([]byte("wrong input,expect start with: \"say:\"\n"))
			return
		}
		body := line[4:]
		conn.Write([]byte("reply:"))
		conn.Write(body)
		conn.Write([]byte("\n"))
	}
}

func (p *Repeater) Close() error {
	atomic.StoreInt32(&p.running, 0)
	return nil
}

var _ servers.Server = &Repeater{}
