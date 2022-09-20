package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var SkillController = &skillController{}

type skillController struct {
}

func (sh *skillController) Router(r *net.Router) {
	g := r.Group("skill")
	g.AddRouter("list", sh.list)
}

func (sh *skillController) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.SkillListRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role, _ := req.Conn.GetProperty("role")
	r := role.(*data.Role)
	skills, err := logic.DefaultSkillService.GetSkills(r.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = skills
	rsp.Body.Msg = rspObj
}
