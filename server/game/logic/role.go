package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"time"
)

var RoleService = &roleService{}

type roleService struct {
}

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, conn net.WSConn) error {
	//根据用户id 查找角色
	role := &data.Role{}
	ok, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "查询角色出错")
	}
	if !ok {
		log.Println("角色不存在", err)
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	rid := role.RId
	rsp.Role = role.ToModel().(model.Role)
	//查询资源
	roleRes := &data.RoleRes{}
	ok, err = db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return common.New(constant.DBError, "查询角色资源出错")
	}
	if !ok {
		//资源不存在  加载初始资源
		roleRes = &data.RoleRes{RId: role.RId,
			Wood:   gameConfig.Base.Role.Wood,
			Iron:   gameConfig.Base.Role.Iron,
			Stone:  gameConfig.Base.Role.Stone,
			Grain:  gameConfig.Base.Role.Grain,
			Gold:   gameConfig.Base.Role.Gold,
			Decree: gameConfig.Base.Role.Decree}

	}
	rsp.RoleRes = roleRes.ToModel().(model.RoleRes)
	rsp.Time = time.Now().UnixNano() / 1e6
	token, err := utils.Award(rid)
	if err != nil {
		log.Println("生成token出错", err)
		return common.New(constant.SessionInvalid, "生成token出错")
	}
	rsp.Token = token
	conn.SetProperty("role", role)
	//初始化角色属性
	if err := DefaultRoleAttrService.TryCreate(rid); err != nil {
		return common.New(constant.DBError, "数据库错误")
	}
	//初始化主城
	err = RoleCityService.InitCity(role)
	if err != nil {
		return common.New(constant.DBError, "数据库错误")
	}
	return nil
}
