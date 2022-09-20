package net

import "sync"

var Mgr = &WsMgr{
	userCache: make(map[int]WSConn),
}

type WsMgr struct {
	uc        sync.RWMutex
	userCache map[int]WSConn
}

func (m *WsMgr) UserLogin(conn WSConn, uid int, token string) {
	m.uc.Lock()
	defer m.uc.Unlock()
	oldConn := m.userCache[uid]
	//当前用户登录着呢
	if oldConn != nil {
		//通知客户端，有用户抢登录了
		if conn != oldConn {
			oldConn.Push("robLogin", nil)
		}
	}
	m.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}
