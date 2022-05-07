// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/29

package repeater

import (
	"bytes"

	"github.com/fsgo/hydra"
)

type Protocol struct{}

func (p *Protocol) MinLen() int {
	return 4
}

func (p *Protocol) MaxLen() int {
	return 4
}

func (p *Protocol) MustNot(header []byte) bool {
	return false
}

func (p *Protocol) Is(header []byte) bool {
	return bytes.HasPrefix(header, []byte("say:"))
}

func (p *Protocol) Name() string {
	return "repeater"
}

var _ hydra.Protocol = (*Protocol)(nil)
