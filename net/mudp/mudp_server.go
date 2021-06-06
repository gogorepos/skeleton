package mudp

import (
	"errors"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/util/gconv"
)

const iDefaultServer = "default"

type Server struct {
	conn    *gudp.Conn
	address string
	handler func(*gudp.Conn)
}

var serverMapping = gmap.NewStrAnyMap(true)

func GetServer(name ...interface{}) *Server {
	serverName := iDefaultServer
	if len(name) > 0 && name[0] != "" {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	s := NewServer("", nil)
	serverMapping.Set(serverName, s)
	return s
}

func (s *Server) SetAddress(address string) {
	s.address = address
}

func (s *Server) SetHandler(handler func(*gudp.Conn)) {
	s.handler = handler
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func (s *Server) Run() error {
	if s.handler == nil {
		return errors.New("start running failed: socket handler not defined")
	}
	conn, err := NewMulticastConn(s.address)
	if err != nil {
		return err
	}
	s.conn = gudp.NewConnByNetConn(conn)
	s.handler(s.conn)
	return nil
}

// NewServer 根据 <address> 和处理函数 <handler> 创建并返回一个 *Server
func NewServer(address string, handler func(*gudp.Conn), name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Set(name[0], s)
	}
	return s
}
