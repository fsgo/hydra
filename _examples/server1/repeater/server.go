// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package repeater

import (
	"bufio"
	"io"
	"net"
	"sync/atomic"
)

type Server struct {
	running int32
	Handler func(input []byte, writer io.Writer)
}

func (p *Server) Serve(l net.Listener) error {
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

func (p *Server) Running() bool {
	return atomic.LoadInt32(&p.running) == 1
}

func (p *Server) handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for p.Running() {
		line, _, err := reader.ReadLine()
		if err != nil {
			conn.Write([]byte("ReadLine with error:"))
			conn.Write([]byte(err.Error()))
			return
		}
		conn.Write([]byte("reply:"))
		conn.Write(line)
		conn.Write([]byte("\n"))
	}
}

func (p *Server) Close() error {
	atomic.StoreInt32(&p.running, 0)
	return nil
}
