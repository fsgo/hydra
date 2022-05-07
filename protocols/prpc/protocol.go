// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package prpc

import (
	"github.com/fsgo/hydra"
)

var expect = []byte{'P', 'R', 'P', 'C'}

type Protocol struct{}

func (p *Protocol) MinLen() int {
	return 1
}

func (p *Protocol) MaxLen() int {
	return 4
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
	return header[0] != expect[0]
}

func (p *Protocol) Name() string {
	return "PRPC"
}

var _ hydra.Protocol = (*Protocol)(nil)
