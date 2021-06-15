package ws

import (
	"errors"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
)

const iDefaultServer = "default"

type Server struct {
	Action  string
	Data    interface{}
	ws      *ghttp.WebSocket
	req     *ghttp.Request
	handler func(*ghttp.WebSocket)
}

var serverMapping = gmap.NewStrAnyMap(true)

func (s *Server) Send(data []byte) error {
	return s.ws.WriteMessage(websocket.TextMessage, data)
}

func (s *Server) SetRequest(req *ghttp.Request) {
	s.req = req
}

func (s *Server) SetHandler(handler func(*ghttp.WebSocket)) {
	s.handler = handler
}

func (s *Server) Run() error {
	if s.req == nil {
		return errors.New("start running filed: websocket req not defined")
	}
	if s.handler == nil {
		return errors.New("start running failed: websocket handler not defined")
	}
	s.handler(s.ws)
	return nil
}

func NewServer(req *ghttp.Request, handler func(*ghttp.WebSocket)) *Server {
	ws, _ := req.WebSocket()
	s := &Server{
		ws:      ws,
		req:     req,
		handler: handler,
	}
	return s
}

func BindServer(handler func(*ghttp.WebSocket), pattern string, name ...interface{}) error {
	var err error
	s := g.Server(name...)
	s.BindHandler(pattern, func(r *ghttp.Request) {
		ws := NewServer(r, handler)
		err = ws.Run()
	})
	return err
}

func Send(data []byte, name ...interface{}) {
	serverName := iDefaultServer
	if len(name) > 0 {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		arr := s.(*garray.Array)
		arr.Walk(func(value interface{}) interface{} {
			server := value.(*Server)
			_ = server.Send(data)
			return value
		})
	}
}
