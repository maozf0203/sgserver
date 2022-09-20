package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"time"
)

var InteriorController = &interiorController{}

type interiorController struct {
}

func (i *interiorController) Router(router *net.Router) {
	g := router.Group("interior")
	g.Use(middleware.Log())
	g.AddRouter("openCollect", i.openCollect, middleware.CheckRole())
}
func (i *interiorController) openCollect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.OpenCollectionRsp{}

	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)
	ra := logic.DefaultRoleAttrService.Get(role.RId)
	if ra == nil {
		rsp.Body.Code = constant.DBError
		return
	}

	interval := gameConfig.Base.Role.CollectInterval
	timeLimit := gameConfig.Base.Role.CollectTimesLimit
	rspObj.Limit = timeLimit
	rspObj.CurTimes = ra.CollectTimes
	if ra.LastCollectTime.IsZero() {
		rspObj.NextTime = 0
	} else {
		if ra.CollectTimes >= timeLimit {
			y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
			//东八区time.FixedZone("CST", 8*3600)
			nextTime := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))
			rspObj.NextTime = nextTime.UnixNano() / 1e6
		} else {
			nextTime := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
			rspObj.NextTime = nextTime.UnixNano() / 1e6
		}
	}
}
