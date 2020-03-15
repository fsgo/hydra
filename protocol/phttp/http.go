/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package phttp

import (
	"bytes"
	"net"
	"net/http"

	"github.com/fsgo/hydra/protocol"
)

var headLen protocol.DiscernLengths

// HTTP 协议
type HTTP struct {
	Server *http.Server
}

func (p *HTTP) MustNot(header []byte) bool {
	first := header[0]
	if _, has := methodFirstBytes[first]; !has {
		return true
	}
	return false
}

// HeaderLen 可判断协议的最小长度
func (p *HTTP) HeaderLen() protocol.DiscernLengths {
	return headLen
}

// Name 协议名称
func (p *HTTP) Name() string {
	return "HTTP"
}

// Is 判断是否当前支持的协议
func (p *HTTP) Is(header []byte) bool {
	spaceIdx := bytes.IndexByte(header, ' ')
	if spaceIdx < minHeaderLength {
		return false
	}
	method := string(bytes.ToUpper(header[:spaceIdx]))

	if _, has := methodsMap[method]; has {
		return true
	}
	return false
}

// Serve 对请求进行处理
func (p *HTTP) Serve(l net.Listener) error {
	if p.Server == nil {
		p.Server = &http.Server{}
	}
	return p.Server.Serve(l)
}

// Close 关闭协议服务
func (p *HTTP) Close() error {
	if p.Server == nil {
		return nil
	}
	return p.Server.Close()
}

var _ protocol.Protocol = &HTTP{}
