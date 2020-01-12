/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/29
 */

package protocol

import (
	"time"
)

// Config 协议配置
type Config struct {
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration
}

// SetReadTimeout 设置读超时
func (c *Config) SetReadTimeout(timeout time.Duration) {
	c.ReadTimeout = timeout
}

// SetWriteTimeout 设置写超时
func (c *Config) SetWriteTimeout(timeout time.Duration) {
	c.WriteTimeout = timeout
}

// SetIdleTimeout 设置等待超时
func (c *Config) SetIdleTimeout(timeout time.Duration) {
	c.IdleTimeout = timeout
}
