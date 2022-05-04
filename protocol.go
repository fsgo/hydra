// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package hydra

import (
	"fmt"
)

// Protocol 协议接口定义
type Protocol interface {
	DiscernLengths() DiscernLengths

	// MustNot 通过 header 判断，一定不是该协议
	// 此时 Header 的长度是 DiscernLengths[0]
	// 若不确定，应该返回 false
	MustNot(header []byte) bool

	// Is 通过header判断是否该协议
	// 此时 Header 的长度是 DiscernLengths[1]
	Is(header []byte) bool

	// Name 协议名字
	Name() string
}

// DiscernLengths 协议头判断的长度：{长度1,长度2}
// 长度1 为最小长度，用于做非判断( MustNot )，用于明确的判断不是该协议
type DiscernLengths [2]int

func (dl DiscernLengths) MustValid() {
	if dl[0] <= 0 || dl[1] <= 0 {
		panic(fmt.Sprintf("DiscernLengths=%v is invalid", dl))
	}
}
