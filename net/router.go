package net

import (
	"log"
	"strings"
	"sync"
)

type Handlerfunc func(req *WsMsgReq, rsp *WsMsgRsp)
type MiddlewareFunc func(handlerFunc Handlerfunc) Handlerfunc

type group struct {
	mutex         sync.Mutex
	prefix        string
	handlerMap    map[string]Handlerfunc
	middlewareMap map[string][]MiddlewareFunc
	middlewares   []MiddlewareFunc
}

func (g *group) AddRouter(name string, h Handlerfunc, middlewares ...MiddlewareFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.handlerMap[name] = h
	g.middlewareMap[name] = middlewares
}

func (g *group) Use(middlewares ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:        prefix,
		handlerMap:    make(map[string]Handlerfunc),
		middlewareMap: make(map[string][]MiddlewareFunc),
	}
	r.group = append(r.group, g)
	return g
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h, ok := g.handlerMap[name]
	if !ok {
		h, ok = g.handlerMap["*"]
		if !ok {
			log.Println("路由未定义")
		}
	}
	if ok {
		//执行中间件 先加的 先执行
		for i := 0; i < len(g.middlewares); i++ {
			h = g.middlewares[i](h)
		}
		mm, ok := g.middlewareMap[name]
		if ok {
			for i := 0; i < len(mm); i++ {
				h = mm[i](h)
			}
		}
		h(req, rsp)
	}
}

type Router struct {
	group []*group
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//req.Body.Name 表示路径由(业务组.路由标识)组成
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		} else if g.prefix == "*" {
			g.exec(name, req, rsp)
		}

	}

}
