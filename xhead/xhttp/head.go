/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/28
 */

package xhttp

import (
	"bytes"
	"strings"

	"github.com/fsgo/hydra"
)

var headLen hydra.DiscernLengths

// Head 协议
type Head struct{}

func (p *Head) MustNot(header []byte) bool {
	first := header[0]
	if _, has := methodFirstBytes[first]; !has {
		return true
	}
	return false
}

// HeaderLen 可判断协议的最小长度
func (p *Head) HeaderLen() hydra.DiscernLengths {
	return headLen
}

// Name 协议名称
func (p *Head) Name() string {
	return "HTTP"
}

// Is 判断是否当前支持的协议
func (p *Head) Is(header []byte) bool {
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

var _ hydra.Head = &Head{}
