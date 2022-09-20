package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var WarReportController = &warReportController{}

type warReportController struct {
}

func (w *warReportController) Router(r *net.Router) {
	g := r.Group("war")
	g.AddRouter("report", w.report)
}

func (w *warReportController) report(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.WarReportRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	role, _ := req.Conn.GetProperty("role")
	r := role.(*data.Role)

	wReports, err := logic.DefaultWarService.GetWarReports(r.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = wReports
	rsp.Body.Msg = rspObj
}
