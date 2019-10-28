package net

import (
  "context"
  "fmt"
  "time"
  "net"
  "sync/atomic"
  "github.com/yggdrasil-network/yggdrasil-go/src/yggdrasil"
)

type YggdrasilListener struct {
	listener  *yggdrasil.Listener
	heartBeat time.Duration
	deadline atomic.Value
}

func NewYggdrasilListener(node YggdrasilNode, heartBeat time.Duration) (*YggdrasilListener, error) {
	listener, err := node.Core.ConnListen()
	if err != nil {
		return nil, fmt.Errorf("cannot create new Yggdrasil listener: %v", err)
	}
	return &YggdrasilListener{listener: listener, heartBeat: heartBeat}, nil
}

// AcceptWithContext waits with context for a generic Conn.
func (l *YggdrasilListener) AcceptWithContext(ctx context.Context) (net.Conn, error) {
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				return nil, fmt.Errorf("cannot accept connections: %v", ctx.Err())
			}
			return nil, nil
		default:
		}
		err := l.SetDeadline(time.Now().Add(l.heartBeat))
		if err != nil {
			return nil, fmt.Errorf("cannot accept connections: %v", err)
		}
		rw, err := l.listener.Accept()
		if err != nil {
			if isTemporary(err) {
				continue
			}
			return nil, fmt.Errorf("cannot accept connections: %v", err)
		}
		return rw, nil
	}
}

// Accept waits for a generic Conn.
func (l *YggdrasilListener) Accept() (net.Conn, error) {
	return l.AcceptWithContext(context.Background())
}

// SetDeadline sets deadline for accept operation.
func (l *YggdrasilListener) SetDeadline(t time.Time) error {
	return l.listener.SetDeadline(t)
}

// Close closes the connection.
func (l *YggdrasilListener) Close() error {
	return l.listener.Close()
}

// Addr represents a network end point address.
func (l *YggdrasilListener) Addr() net.Addr {
	return l.listener.Addr()
}
