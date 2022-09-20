package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type server struct {
	addr       string
	router     *Router
	needSecret bool
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) NeedSecret(needSecret bool) {
	s.needSecret = needSecret
}

func (s *server) Router(router *Router) {
	s.router = router
}

// 启动服务
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket服务链接出错")
	}
	wsServer := NewWsServer(wsConn, s.needSecret)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}
