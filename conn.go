// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/5/4

package hydra

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

type picker interface {
	Pick(n int) ([]byte, error)
}

func newReferee(h Protocol, lis Listener) *referee {
	return &referee{
		protocol: h,
		listener: lis,
	}
}

type referees []*referee

func (rs referees) MustCheck() {
	names := make(map[string]bool, len(rs))
	for i := 0; i < len(rs); i++ {
		name := rs[i].protocol.Name()
		if names[name] {
			panic(fmt.Sprintf("cannot bind protocol %s twice", name))
		}
		names[name] = true
	}
}

func (rs referees) Sort() {
	sort.Slice(rs, func(i, j int) bool {
		a := rs[i]
		b := rs[j]
		return a.protocol.MaxLen() > b.protocol.MaxLen()
	})
}

func (rs referees) Close() {
	for i := 0; i < len(rs); i++ {
		_ = rs[i].listener.Close()
	}
}

type referee struct {
	listener Listener
	protocol Protocol
}

func (s *referee) Is(conn picker) (is bool, err error) {
	// 先用少量字符做明确的非判断
	header, errHeader := conn.Pick(s.protocol.MinLen())
	if errHeader != nil {
		return false, errHeader
	}

	if s.protocol.MustNot(header) {
		return false, nil
	}

	// 再验证更多字符
	header, errHeader = conn.Pick(s.protocol.MaxLen())
	if errHeader != nil {
		return false, errHeader
	}

	if s.protocol.Is(header) {
		return true, nil
	}
	return false, nil
}

// NewListener 创建一个指定  chan size 的 Listener
func NewListener(size int) Listener {
	return &listener{
		Connects: make(chan net.Conn, size),
	}
}

// Listener 具备接收 Conn 的 Listener
type Listener interface {
	net.Listener
	Dispatch(conn net.Conn)
	SetLocalAddr(addr net.Addr)
}

var _ Listener = (*listener)(nil)

type listener struct {
	localAddr        net.Addr
	Connects         chan net.Conn
	mux              sync.Mutex
	dl               time.Time
	closed           bool
	errAcceptTimeout *net.OpError
}

func (l *listener) SetLocalAddr(addr net.Addr) {
	l.localAddr = addr
	l.initErrors()
}

func (l *listener) initErrors() {
	if l.errAcceptTimeout == nil {
		l.errAcceptTimeout = &net.OpError{
			Op:   "accept",
			Net:  l.localAddr.Network(),
			Addr: l.localAddr,
			Err: &netErr{
				timeout:   true,
				temporary: true,
				msg:       "i/o timeout",
			},
		}
	}
}

type canSetDeadline interface {
	SetDeadline(t time.Time) error
}

func (l *listener) SetDeadline(t time.Time) error {
	l.mux.Lock()
	l.dl = t
	l.mux.Unlock()
	return nil
}

var errListenerClosed = errors.New("listener closed")

func (l *listener) Accept() (net.Conn, error) {
	l.initErrors()

	l.mux.Lock()
	dl := l.dl
	closed := l.closed
	l.mux.Unlock()

	if closed {
		return nil, errListenerClosed
	}

	if dl.IsZero() {
		if c, ok := <-l.Connects; ok {
			return c, nil
		}
		return nil, errListenerClosed
	}

	timeout := time.Until(dl)
	if timeout <= 0 {
		return nil, l.errAcceptTimeout
	}

	tm := time.NewTimer(timeout)
	defer tm.Stop()

	select {
	case c := <-l.Connects:
		return c, nil
	case <-tm.C:
		return nil, l.errAcceptTimeout
	}
}

func (l *listener) Close() error {
	l.mux.Lock()
	l.mux.Unlock()
	if l.closed {
		return nil
	}
	l.closed = true
	close(l.Connects)
	return nil
}

func (l *listener) Addr() net.Addr {
	return l.localAddr
}

func (l *listener) Dispatch(conn net.Conn) {
	l.mux.Lock()
	closed := l.closed
	l.mux.Unlock()
	if closed {
		_ = conn.Close()
		return
	}
	l.Connects <- conn
}

func newConn(raw net.Conn) *pickConn {
	return &pickConn{
		Conn: raw,
		bio:  bufio.NewReader(raw),
	}
}

type pickConn struct {
	net.Conn
	bio *bufio.Reader
}

// Pick 获取请求头，用于判断协议
func (c *pickConn) Pick(expectLen int) ([]byte, error) {
	return c.bio.Peek(expectLen)
}

// Read 读取内容
func (c *pickConn) Read(b []byte) (int, error) {
	return c.bio.Read(b)
}

func (c *pickConn) Close() error {
	return c.Conn.Close()
}

var _ error = (*netErr)(nil)

type netErr struct {
	msg       string
	timeout   bool
	temporary bool
}

func (ne *netErr) Error() string {
	return ne.msg
}

func (ne *netErr) Timeout() bool {
	return ne.timeout
}
func (ne *netErr) Temporary() bool {
	return ne.temporary
}
