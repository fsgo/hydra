// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package hydra

// Protocol 协议接口定义
type Protocol interface {
	// Name 协议名字
	Name() string

	// MinLen 为最小长度，用于做非判断( MustNot )
	MinLen() int

	// MustNot 通过 MinLen 长度的 header 判断，一定不是该协议
	// 	当确定不是该协议的时候，返回 true
	// 	若不确定，应该返回 false
	MustNot(header []byte) bool

	// MaxLen 最大长度，用于明确的判断不是该协议
	// 应确保 MaxLen >= MinLen > 0
	MaxLen() int

	// Is 通过 MaxLen 长度的 header 判断是否该协议
	Is(header []byte) bool
}
