package game

import (
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/game/controller"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/gameConfig/general"
	"sgserver/server/game/logic"
)

var Router = &net.Router{}

func Init() {
	//初始化数据库
	db.TestDB()
	//加载基础配置
	gameConfig.Base.Load()
	//加载地图资源
	gameConfig.MapBuildConf.Load()
	//加载城池设施配置
	gameConfig.FacilityConf.Load()
	//加载地图配置
	gameConfig.MapRes.Load()
	//初始化武将配置
	general.General.Load()
	//加载技能配置
	gameConfig.Skill.Load()
	//加载所有建筑信息
	logic.DefaultRoleBuildService.Load()
	//加载城池信息
	logic.RoleCityService.Load()
	//加载所有角色属性
	logic.DefaultRoleAttrService.Load()
	//初始化路由
	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
	controller.DefaultGeneralController.Router(Router)
	controller.DefaultArmyController.Router(Router)
	controller.WarReportController.Router(Router)
	controller.SkillController.Router(Router)
	controller.InteriorController.Router(Router)
}
