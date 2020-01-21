/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package servers

import (
	"net"
)

// Server 协议接口定义
type Server interface {
	HeaderLen() int

	Is(header []byte) bool

	Name() string

	Serve(l net.Listener) error
	Close() error
}
