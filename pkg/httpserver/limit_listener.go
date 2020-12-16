package httpserver

import (
	"net"
	"sync"
)

type limitedListener struct {
	net.Listener
	sem chan struct{}
}

type limitedListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

func NewLimitedListener(l net.Listener, maxRequests uint) net.Listener {
	return &limitedListener{l, make(chan struct{}, maxRequests)}
}

func (l *limitedListener) acquire() { l.sem <- struct{}{} }
func (l *limitedListener) release() { <-l.sem }

func (l *limitedListener) Accept() (net.Conn, error) {
	l.acquire()

	c, err := l.Listener.Accept()
	if err != nil {
		l.release()

		return nil, err
	}

	return &limitedListenerConn{
		Conn:    c,
		release: l.release,
	}, nil
}

func (l *limitedListenerConn) Close() error {
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)

	return err
}
