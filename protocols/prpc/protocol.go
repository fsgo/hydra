// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package prpc

import (
	"github.com/fsgo/hydra"
)

var headLen = hydra.DiscernLengths{1, 4}

var expect = []byte{'P', 'R', 'P', 'C'}

type Protocol struct {
}

func (p *Protocol) Is(header []byte) bool {
	for i := 0; i < len(expect); i++ {
		if header[i] != expect[i] {
			return false
		}
	}
	return true
}

func (p *Protocol) MustNot(header []byte) bool {
	first := header[0]
	return first != header[0]
}

func (p *Protocol) Name() string {
	// TODO implement me
	panic("implement me")
}

func (p *Protocol) DiscernLengths() hydra.DiscernLengths {
	return headLen
}

var _ hydra.Protocol = &Protocol{}
