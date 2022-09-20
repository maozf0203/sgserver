package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var DefaultGeneralController = &generalController{}

type generalController struct {
}

func (gh *generalController) Router(r *net.Router) {
	g := r.Group("general")
	g.AddRouter("myGenerals", gh.myGenerals)
}

func (gh *generalController) myGenerals(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.MyGeneralRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role, _ := req.Conn.GetProperty("role")
	r := role.(*data.Role)
	gs, err := logic.DefaultGeneralService.GetGenerals(r.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Generals = gs
	rsp.Body.Msg = rspObj
}
