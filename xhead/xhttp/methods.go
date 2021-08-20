// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/15

package xhttp

import (
	"strings"
)

var methods = []string{
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

// header中空格出现的最小位置
var minHeaderLength = -1

var methodsMap = map[string]bool{}

// 初始化
func Init() {
	headLen = [2]int{1, 1}

	for _, methodUpper := range methods {

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
	Init()
}
