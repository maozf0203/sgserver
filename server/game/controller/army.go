package controller

import (
	"github.com/mitchellh/mapstructure"
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var DefaultArmyController = &armyController{}

type armyController struct {
}

func (gh *armyController) Router(r *net.Router) {
	g := r.Group("army")
	g.AddRouter("myList", gh.myList)
}

func (ah *armyController) myList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ArmyListReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.ArmyListRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role, _ := req.Conn.GetProperty("role")
	r := role.(*data.Role)
	arms, err := logic.DefaultArmyService.GetArmysByCity(r.RId, reqObj.CityId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Armys = arms
	rspObj.CityId = reqObj.CityId
	rsp.Body.Msg = rspObj
}
