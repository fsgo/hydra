// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/29

package repeater

import (
	"bytes"

	"github.com/fsgo/hydra"
)

type Head struct{}

func (p *Head) DiscernLengths() hydra.DiscernLengths {
	return [2]int{4, 4}
}

func (p *Head) MustNot(header []byte) bool {
	return false
}

func (p *Head) Is(header []byte) bool {
	return bytes.HasPrefix(header, []byte("say:"))
}

func (p *Head) Name() string {
	return "repeater"
}

var _ hydra.Protocol = &Head{}
