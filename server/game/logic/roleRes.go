package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var RoleResService = &roleResService{}

type roleResService struct {
}

func (r *roleResService) GetRoleRes(rid int) (model.RoleRes, error) {
	roleRes := &data.RoleRes{}
	ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return model.RoleRes{}, common.New(constant.DBError, "查询角色资源出错")
	}
	if ok {
		return roleRes.ToModel().(model.RoleRes), nil
	}
	return model.RoleRes{}, nil
}
