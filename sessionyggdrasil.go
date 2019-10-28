package coap

import (
	"bytes"
	"context"
	"fmt"
	"net"

	coapNet "github.com/Fnux/go-coap/net"
)

type sessionYggdrasil struct {
	sessionBase
	connection *coapNet.Conn
}

// newSessionYggdrasil create new session for Yggdrasil connection
func newSessionYggdrasil(connection *coapNet.Conn, srv *Server) (networkSession, error) {
	BlockWiseTransfer := false
	BlockWiseTransferSzx := BlockWiseSzx1024
	if srv.BlockWiseTransfer != nil {
		BlockWiseTransfer = *srv.BlockWiseTransfer
	}
	if srv.BlockWiseTransferSzx != nil {
		BlockWiseTransferSzx = *srv.BlockWiseTransferSzx
	}

	if BlockWiseTransfer && BlockWiseTransferSzx == BlockWiseSzxBERT {
		return nil, ErrInvalidBlockWiseSzx
	}

	s := sessionYggdrasil{
		connection: connection,
		sessionBase: sessionBase{
			srv:                  srv,
			handler:              &TokenHandler{tokenHandlers: make(map[[MaxTokenSize]byte]HandlerFunc)},
			blockWiseTransfer:    BlockWiseTransfer,
			blockWiseTransferSzx: uint32(BlockWiseTransferSzx),
			mapPairs:             make(map[[MaxTokenSize]byte]map[uint16](*sessionResp)),
		},
	}

	return &s, nil
}

// LocalAddr implements the networkSession.LocalAddr method.
func (s *sessionYggdrasil) LocalAddr() net.Addr {
	return s.connection.LocalAddr()
}

// RemoteAddr implements the networkSession.RemoteAddr method.
func (s *sessionYggdrasil) RemoteAddr() net.Addr {
	return s.connection.RemoteAddr()
}

// BlockWiseTransferEnabled
func (s *sessionYggdrasil) blockWiseEnabled() bool {
	return s.blockWiseTransfer
}

func (s *sessionYggdrasil) blockWiseIsValid(szx BlockWiseSzx) bool {
	return szx <= BlockWiseSzx1024
}

// Ping send ping over udp(unicast) and wait for response.
func (s *sessionYggdrasil) PingWithContext(ctx context.Context) error {
  // TODO: check if relevant for Yggdrasil.
	//provoking to get a reset message - "CoAP ping" in RFC-7252
	//https://tools.ietf.org/html/rfc7252#section-4.2
	//https://tools.ietf.org/html/rfc7252#section-4.3
	//https://tools.ietf.org/html/rfc7252#section-1.2 "Reset Message"
	// BUG of iotivity: https://jira.iotivity.org/browse/IOT-3149
	req := s.NewMessage(MessageParams{
		Type:      Confirmable,
		Code:      Empty,
		MessageID: GenerateMessageID(),
	})
	resp, err := s.ExchangeWithContext(ctx, req)
	if err != nil {
		return err
	}
	if resp.Type() == Reset {
		return nil
	}
	return ErrInvalidResponse
}

func (s *sessionYggdrasil) closeWithError(err error) error {
	if s.connection != nil {
		c := ClientConn{commander: &ClientCommander{s}}
		s.srv.NotifySessionEndFunc(&c, err)
		e := s.connection.Close()
		//s.connection = nil
		if e == nil {
			e = err
		}
		return e
	}
	return err
}

// Close implements the networkSession.Close method
func (s *sessionYggdrasil) Close() error {
	return s.closeWithError(nil)
}

// NewMessage Create message for response
func (s *sessionYggdrasil) NewMessage(p MessageParams) Message {
	return NewDgramMessage(p)
}

func (s *sessionYggdrasil) IsTCP() bool {
	return false
}

func (s *sessionYggdrasil) ExchangeWithContext(ctx context.Context, req Message) (Message, error) {
	return s.exchangeWithContext(ctx, req, s.WriteMsgWithContext)
}

// Write implements the networkSession.Write method.
func (s *sessionYggdrasil) WriteMsgWithContext(ctx context.Context, req Message) error {
	buffer := bytes.NewBuffer(make([]byte, 0, 1500))
	err := req.MarshalBinary(buffer)
	if err != nil {
		return fmt.Errorf("cannot write msg to tcp connection %v", err)
	}
	return s.connection.WriteWithContext(ctx, buffer.Bytes())
}

func (s *sessionYggdrasil) sendPong(w ResponseWriter, r *Request) error {
	resp := r.Client.NewMessage(MessageParams{
		Type:      Reset,
		Code:      Empty,
		MessageID: r.Msg.MessageID(),
	})
	return w.WriteMsgWithContext(r.Ctx, resp)
}

func (s *sessionYggdrasil) handleSignals(w ResponseWriter, r *Request) bool {
	switch r.Msg.Code() {
	// handle of udp ping
	case Empty:
		if r.Msg.Type() == Confirmable && r.Msg.AllOptions().Len() == 0 && (r.Msg.Payload() == nil || len(r.Msg.Payload()) == 0) {
			s.sendPong(w, r)
			return true
		}
	}
	return false
}
