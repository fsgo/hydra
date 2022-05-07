// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/28

package xhttp

import (
	"bytes"
	"strings"

	"github.com/fsgo/hydra"
)

var headLen [2]int

// Protocol 协议
type Protocol struct{}

func (p *Protocol) MinLen() int {
	return headLen[0]
}

func (p *Protocol) MaxLen() int {
	return headLen[1]
}

// MustNot 通过首字母，快速进行非判断
func (p *Protocol) MustNot(header []byte) bool {
	first := header[0]
	if _, has := methodFirstBytes[first]; !has {
		return true
	}
	return false
}

// Name 协议名称
func (p *Protocol) Name() string {
	return "HTTP"
}

// Is 判断是否当前支持的协议
func (p *Protocol) Is(header []byte) bool {
	spaceIdx := bytes.IndexByte(header, ' ')
	if spaceIdx < minHeaderLength {
		return false
	}

	method := string(header[:spaceIdx])
	if _, has := methodsMap[method]; has {
		return true
	}

	method = strings.ToUpper(method)
	if _, has := methodsMap[method]; has {
		return true
	}

	return false
}

var _ hydra.Protocol = (*Protocol)(nil)
