// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package protocols

import (
	"github.com/fsgo/hydra"
	"github.com/fsgo/hydra/protocols/prpc"
	"github.com/fsgo/hydra/protocols/xhttp"
)

func HTTP() hydra.Protocol {
	return &xhttp.Protocol{}
}

func PRPC() hydra.Protocol {
	return &prpc.Protocol{}
}
