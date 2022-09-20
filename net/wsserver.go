package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/utils"
	"sync"
	"time"
)

// websocket服务
type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
	needSecret   bool
}

var cid int64

func NewWsServer(wsConn *websocket.Conn, needSecret bool) *wsServer {
	s := &wsServer{
		wsConn:     wsConn,
		outChan:    make(chan *WsMsgRsp, 1000),
		property:   make(map[string]interface{}),
		Seq:        0,
		needSecret: needSecret,
	}
	cid++
	s.SetProperty("cid", cid)
	return s
}

func (w *wsServer) Router(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}
func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{
		Body: &RspBody{
			Name: name,
			Msg:  data,
			Seq:  0,
		},
	}
	w.outChan <- rsp
}

// 通道一旦建立，那么 收发消息就要一直监听
func (w *wsServer) Start() {
	//启动读取数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			w.Write(msg.Body)
		}
	}
}

func (w *wsServer) Write(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("服务端写数据:" + string(data))
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		//加密
		key := secretKey.(string)
		//数据加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	//压缩
	if data, err = utils.Zip(data); err == nil {
		w.wsConn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func (w *wsServer) readMsgLoop() {
	//先读到客户端发送的消息，然后进行处理，然后回复消息
	//经过路由 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("接收消息出现错误:", err)
			break
		}
		//收到消息 解析消息 前端发过来的消息就是json格式
		//1.data 解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式", err)
			continue
		}
		//2.前端的消息是经过加密的消息，需要解密
		if w.needSecret {
			secretKey, err := w.GetProperty("secretKey")
			if err == nil {
				//有加密
				key := secretKey.(string)
				//客户端传过来的数据是加密的 需要解密
				d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
				if err != nil {
					log.Println("数据格式有误，解密失败:", err)
					//出错后 发起握手
					w.Handshake()
				} else {
					data = d
				}
			}
		}
		fmt.Println("服务端接收消息:" + string(data))
		//3.格式化消息
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("服务端json数据格式有误，非法格式:", err)
		} else {
			//获取到前端传递的数据了，进行业务处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			if req.Body.Name == HeartbeatMsg {
				//回心跳消息
				h := &Heartbeat{}
				mapstructure.Decode(body.Msg, h)
				h.STime = time.Now().UnixNano() / 1e6
				rsp.Body.Msg = h
			} else {
				if w.router != nil {
					w.router.Run(req, rsp)
				}
			}
			w.outChan <- rsp
		}
	}
	w.Close()
}

func (w *wsServer) Close() {
	w.wsConn.Close()
}

const HandshakeMsg = "handshake"

// 当游戏客户端 发送请求的时候 会先进行握手协议
// 后端会发送对应的加密key给客户端
// 客户端在发送数据的时候 就会使用此key进行加密处理
func (w *wsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string)
	} else {
		secretKey = utils.RandSeq(16)
	}
	handshake := &Handshake{Key: secretKey}

	body := &RspBody{Name: HandshakeMsg, Msg: handshake}

	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data, err := utils.Zip(data); err == nil {
			w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
