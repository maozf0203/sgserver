package logic

import (
	"encoding/json"
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sync"
)

var DefaultRoleAttrService = &RoleAttrService{
	attrs: make(map[int]*data.RoleAttribute),
}

type RoleAttrService struct {
	mutex sync.RWMutex
	attrs map[int]*data.RoleAttribute
}

func (r *RoleAttrService) Load() {
	ras := make(map[int]*data.RoleAttribute)
	err := db.Engine.Table(new(data.RoleAttribute)).Find(&ras)
	if err != nil {
		log.Println("roleAttribute Load err", err)
	}
	for _, v := range ras {
		r.attrs[v.RId] = v
	}
}

func (r *RoleAttrService) TryCreate(rid int) error {
	rr := &data.RoleAttribute{}
	ok, err := db.Engine.Table(rr).Where("rid=?", rid).Get(rr)
	if err != nil {
		log.Println("角色属性查询出错", err)
		return err
	}
	if ok {
		r.mutex.Lock()
		r.attrs[rid] = rr
		r.mutex.Unlock()
		return nil
	} else {
		//查询没有 进行初始化创建
		rr.RId = rid
		rr.ParentId = 0
		rr.UnionId = 0
		rr.PosTags = ""
		_, err := db.Engine.Table(rr).Insert(rr)
		if err != nil {
			log.Println("角色属性插入出错", err)
			return err
		}
		r.mutex.Lock()
		r.attrs[rid] = rr
		r.mutex.Unlock()
	}
	return nil
}

func (r *RoleAttrService) GetPosTags(rid int) []model.PosTag {
	r.mutex.RLock()
	rr, ok := r.attrs[rid]
	r.mutex.RUnlock()
	postTags := make([]model.PosTag, 0)
	if ok {
		err := json.Unmarshal([]byte(rr.PosTags), &postTags)
		if err != nil {
			log.Println("标记格式错误", err)
		}
	}
	return postTags
}

func (r *RoleAttrService) Get(rid int) *data.RoleAttribute {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ra, ok := r.attrs[rid]
	if ok {
		return ra
	}
	return nil
}
