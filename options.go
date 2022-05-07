// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package hydra

import (
	"time"
)

type Options struct {
	ListerChanSize int
	AcceptTimeout  time.Duration

	OnAcceptError func(err error)
}

func (o *Options) GetListerChanSize() int {
	if o.ListerChanSize > 0 {
		return o.ListerChanSize
	}
	return 1024
}

func (o *Options) GetAcceptTimeout() time.Duration {
	if o.AcceptTimeout > 0 {
		return o.AcceptTimeout
	}
	return time.Second
}

func (o *Options) invokeOnAcceptError(err error) {
	if o.OnAcceptError == nil {
		return
	}
	o.OnAcceptError(err)
}

var optionsEmpty = &Options{}
