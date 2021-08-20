// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/29

package internal

import (
	"fmt"
)

// Head 协议接口定义
type Head interface {
	HeaderLen() DiscernLengths

	// 通过header判断是否该协议
	Is(header []byte) bool

	// 通过header判断，一定不是该协议
	MustNot(header []byte) bool

	// 协议名字
	Name() string
}

// 协议头判断的长度：{长度1,长度2}
// 长度1 为最小长度，用于做非判断(MustNot)，用于明确的判断不是该协议
type DiscernLengths [2]int

func (dl DiscernLengths) MustValid() {
	if dl[0] <= 0 || dl[1] <= 0 {
		panic(fmt.Sprintf("DiscernLengths=%v is invalid", dl))
	}
}
