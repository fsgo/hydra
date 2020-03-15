/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/3/15
 */

package phttp

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

var methodFirstBytes = map[byte]int{}

// header中空格出现的最小位置
var minHeaderLength = -1

var methodsMap = map[string]int{}

func init() {
	headLen = [2]int{1, 1}
	for _, method := range methods {
		length := len(method)
		if headLen[1] == 0 || headLen[1] < length {
			headLen[1] = length
		}

		first := method[0]
		methodFirstBytes[first] = 1

		methodsMap[method] = 1

		if minHeaderLength < 0 || minHeaderLength > length {
			minHeaderLength = length
		}
	}
}
