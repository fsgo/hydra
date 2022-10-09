// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/15

package xhttp

import (
	"strings"
)

// Methods 支持的 HTTP methods 列表
//
//	若有自定义的 Method 需要支持，请直接修改该变量
//	修改完成后调用 Prepare 方法重新构建
var Methods = []string{
	"GET",
	"POST",
	"PUT",

	"DELETE",
	"HEAD",

	"CONNECT",
	"OPTIONS",
	"PATCH",
	"TRACE",
}

var methodFirstBytes = map[byte]bool{}

// header 中空格出现的最小位置
var minHeaderLength = -1

var methodsMap = map[string]bool{}

// Prepare 构建协议解析所需要的缓存数据
// 默认情况下已在 init 方法里调用
func Prepare() {
	methodsMap = map[string]bool{}
	headLen = [2]int{1, 1}
	minHeaderLength = -1

	for _, methodUpper := range Methods {
		{
			// 获取method 的长度区间
			length := len(methodUpper)
			if headLen[1] == 0 || headLen[1] < length {
				headLen[1] = length
			}
			if minHeaderLength < 0 || minHeaderLength > length {
				minHeaderLength = length
			}
		}

		methodLower := strings.ToLower(methodUpper) // 小写

		{
			// 构建首字母的map表，用于not 判断
			methodFirstBytes[methodUpper[0]] = true // 大小首字母

			methodFirstBytes[methodLower[0]] = true
		}

		{
			// 构建请求方法的 map 表
			methodsMap[methodUpper] = true
			methodsMap[methodLower] = true
		}
	}
}

func init() {
	Prepare()
}
