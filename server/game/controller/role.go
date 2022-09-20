package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"time"
)

var DefaultRoleController = RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.Use(middleware.Log())
	g.AddRouter("enterServer", r.enterServer)
	g.AddRouter("myProperty", r.myProperty, middleware.CheckRole())
	g.AddRouter("posTagList", r.posTagList)
	g.AddRouter("create", r.create)

}

func (r *RoleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//进入游戏的逻辑
	//Session 需要验证是否合法 合法的情况下 可以取出登录用的用户id
	//根据用户id 去查询对应的游戏角色，如果有 就继续 如果没有 提示无角色
	//根据角色id 查询角色拥有的资源，如果有资源就返回，如果没有初始化资源
	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		rsp.Body.Code = constant.InvalidParam
		return
	}
	token := reqObj.Session
	_, claim, err := utils.ParseToken(token)
	if err != nil {
		log.Println("session无效", err)
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	//用户id
	uid := claim.Uid
	err = logic.RoleService.EnterServer(uid, rspObj, req.Conn)
	if err != nil {
		rsp.Body.Msg = rspObj
		rspObj.Time = time.Now().UnixNano() / 1e6
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (rh *RoleController) myProperty(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//分别根据角色id去查询军队 资源  建筑  城池  武将
	reqObj := &model.MyRolePropertyReq{}
	rspObj := &model.MyRolePropertyRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	//城池
	var err error
	rspObj.Citys, err = logic.RoleCityService.GetCitys(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//建筑
	rspObj.MRBuilds, err = logic.DefaultRoleBuildService.GetBuilds(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//资源
	rspObj.RoleRes, err = logic.RoleResService.GetRoleRes(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//武将
	rspObj.Generals, err = logic.DefaultGeneralService.GetGenerals(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//军队
	rspObj.Armys, err = logic.DefaultArmyService.GetArmys(role.RId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (rh *RoleController) posTagList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.PosTagListRsp{}
	role, err := req.Conn.GetProperty("role")
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		rsp.Body.Code = constant.InvalidParam
		return
	}
	r := role.(*data.Role)
	rspObj.PosTags = logic.DefaultRoleAttrService.GetPosTags(r.RId)
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (r *RoleController) create(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.CreateRoleReq{}
	rspObj := &model.CreateRoleRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role := &data.Role{}
	ok, err := db.Engine.Where("uid=?", reqObj.UId).Get(role)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	if ok {
		rsp.Body.Code = constant.RoleAlreadyCreate
		return
	}
	role.UId = reqObj.UId
	role.Sex = reqObj.Sex
	role.NickName = reqObj.NickName
	role.Balance = 0
	role.HeadId = reqObj.HeadId
	role.CreatedAt = time.Now()
	role.LoginTime = time.Now()
	_, err = db.Engine.InsertOne(role)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	rspObj.Role = role.ToModel().(model.Role)
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
