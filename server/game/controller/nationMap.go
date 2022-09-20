package controller

import (
	"github.com/mitchellh/mapstructure"
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var DefaultNationMapController = NationMapController{}

type NationMapController struct {
}

func (n *NationMapController) Router(router *net.Router) {
	g := router.Group("nationMap")
	g.AddRouter("config", n.config)
	g.AddRouter("scanBlock", n.scanBlock, middleware.CheckRole())

}

func (n *NationMapController) config(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ConfigReq{}
	rspObj := &model.ConfigRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	m := gameConfig.MapBuildConf.Cfg
	rspObj.Confs = make([]model.Conf, len(m))
	for index, v := range m {
		rspObj.Confs[index].Type = v.Type
		rspObj.Confs[index].Name = v.Name
		rspObj.Confs[index].Level = v.Level
		rspObj.Confs[index].Defender = v.Defender
		rspObj.Confs[index].Durable = v.Durable
		rspObj.Confs[index].Grain = v.Grain
		rspObj.Confs[index].Iron = v.Iron
		rspObj.Confs[index].Stone = v.Stone
		rspObj.Confs[index].Wood = v.Wood
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
}

func (n *NationMapController) scanBlock(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ScanBlockReq{}
	rspObj := &model.ScanRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK
	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)
	//扫描角色建筑
	mrb := logic.DefaultRoleBuildService.ScanBlock(reqObj)
	rspObj.MRBuilds = mrb
	//扫描角色城池
	mrc := logic.RoleCityService.ScanBlock(reqObj)
	rspObj.MCBuilds = mrc
	//扫描玩家军队
	armys := logic.DefaultArmyService.ScanBlock(role.RId, reqObj)
	rspObj.Armys = armys
	rsp.Body.Msg = rspObj

}
