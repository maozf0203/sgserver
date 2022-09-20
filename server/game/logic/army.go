package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/global"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"sync"
)

var DefaultArmyService = &ArmyService{
	passByPosArmys: make(map[int]map[int]*data.Army),
}

type ArmyService struct {
	passBy         sync.RWMutex
	passByPosArmys map[int]map[int]*data.Army //玩家路过位置的军队 key:posId,armyId
}

func (g *ArmyService) GetArmys(rid int) ([]model.Army, error) {
	mrs := make([]data.Army, 0)
	mr := &data.Army{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, common.New(constant.DBError, "军队查询出错")
	}
	modelMrs := make([]model.Army, len(mrs))
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.Army))
	}
	return modelMrs, nil
}

func (g *ArmyService) GetArmysByCity(rid, cid int) ([]model.Army, error) {
	mrs := make([]data.Army, 0)
	mr := &data.Army{}
	err := db.Engine.Table(mr).Where("rid=? and cityId=?", rid, cid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, common.New(constant.DBError, "军队查询出错")
	}
	modelMrs := make([]model.Army, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.Army))
	}
	return modelMrs, nil
}

func (a *ArmyService) ScanBlock(roleId int, req *model.ScanBlockReq) []model.Army {
	x := req.X
	y := req.Y
	length := req.Length
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return nil
	}

	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)
	out := make([]model.Army, 0)

	a.passBy.RLock()
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {

			posId := global.ToPosition(i, j)
			armys, ok := a.passByPosArmys[posId]
			if ok {
				//是否在视野范围内
				is := armyIsInView(roleId, i, j)
				if is == false {
					continue
				}
				for _, army := range armys {
					out = append(out, army.ToModel().(model.Army))
				}
			}
		}
	}
	a.passBy.RUnlock()
	return out
}

func armyIsInView(rid, x, y int) bool {
	//简单点 先设为true
	return true
}
