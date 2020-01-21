/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package httpserver

import (
	"bytes"
	"net"
	"net/http"

	"github.com/fsgo/hydra/protocol"
)

// Protocol HTTP 协议
type Protocol struct {
	Server *http.Server
}

// HeaderLen 可判断协议的最小长度
func (p *Protocol) HeaderLen() int {
	return 7
}

// Name 协议名称
func (p *Protocol) Name() string {
	return "HTTP"
}

var methods = [][]byte{
	[]byte("GET "),
	[]byte("POST "),
	[]byte("PUT "),

	[]byte("DELETE "),
	[]byte("HEAD "),

	[]byte("CONNECT "),
	[]byte("OPTIONS "),
	[]byte("PATCH "),
	[]byte("TRACE "),
}

// Is 判断是否当前支持的协议
func (p *Protocol) Is(header []byte) bool {
	for _, method := range methods {
		if bytes.HasPrefix(header, method) {
			return true
		}
	}
	firstLine := bytes.SplitN(header, []byte("\r\n"), 2)[0]
	if bytes.Contains(firstLine, []byte("HTTP/")) {
		return true
	}
	return false
}

// Serve 对请求进行处理
func (p *Protocol) Serve(l net.Listener) error {
	if p.Server == nil {
		p.Server = &http.Server{}
	}
	return p.Server.Serve(l)
}

// Close 关闭协议服务
func (p *Protocol) Close() error {
	if p.Server == nil {
		return nil
	}
	return p.Server.Close()
}

var _ protocol.Protocol = &Protocol{}
