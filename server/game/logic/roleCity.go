package logic

import (
	"log"
	"math/rand"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/global"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"sync"
	"time"
)

var RoleCityService = &roleCityService{
	posCity:  make(map[int]*data.MapRoleCity),
	roleCity: make(map[int][]*data.MapRoleCity),
}

type roleCityService struct {
	mutex    sync.RWMutex
	posCity  map[int]*data.MapRoleCity
	roleCity map[int][]*data.MapRoleCity
}

func (r *roleCityService) InitCity(role *data.Role) error {
	rc := &data.MapRoleCity{}
	ok, err := db.Engine.Table(rc).Where("rid=?", role.RId).Get(rc)
	if err != nil {
		log.Println("查询角色城市出错", err)
		return common.New(constant.DBError, "查询角色城市出错")
	}
	if !ok {
		for {
			//没有城池 初始化 条件系统城市5格内不能有玩家城池
			x := rand.Intn(global.MapWith)
			y := rand.Intn(global.MapHeight)
			//判断是否符合创建条件（系统城池五格之内，不能有玩家城池）
			if r.IsCanBuild(x, y) {
				//建的肯定是主城
				rc.RId = role.RId
				rc.Y = y
				rc.X = x
				rc.CreatedAt = time.Now()
				rc.Name = role.NickName
				rc.CurDurable = gameConfig.Base.City.Durable
				rc.IsMain = 1
				_, err := db.Engine.Table(rc).Insert(rc)
				if err != nil {
					log.Println("插入玩家城市出错", err)
					return common.New(constant.DBError, "插入玩家城市出错")
				}
				posId := global.ToPosition(rc.X, rc.Y)
				r.posCity[posId] = rc
				if _, ok := r.roleCity[role.RId]; !ok {
					r.roleCity[role.RId] = make([]*data.MapRoleCity, 0)
				} else {
					r.roleCity[role.RId] = append(r.roleCity[role.RId], rc)
				}
				//生成城市设施
				if err := CityFacilityService.TryCreate(rc.CityId, role.RId); err != nil {
					log.Println("插入城池设施出错", err)
					return common.New(err.(*common.MyError).Code(), err.Error())
				}
				break
			}
		}
	}
	return nil
}

func (r *roleCityService) IsCanBuild(x int, y int) bool {
	confs := gameConfig.MapRes.Confs
	pIndex := global.ToPosition(x, y)
	_, ok := confs[pIndex]
	if !ok {
		return false
	}

	//城池 1范围内 不能超过边界
	if x+1 >= global.MapWith || y+1 >= global.MapHeight || y-1 < 0 || x-1 < 0 {
		return false
	}
	sysBuild := gameConfig.MapRes.SysBuild
	//系统城池的5格内 不能创建玩家城池
	for _, v := range sysBuild {
		if v.Type == gameConfig.MapBuildSysCity {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				y >= v.Y-5 &&
				y <= v.Y+5 {
				return false
			}
		}
	}

	//玩家城池的5格内 也不能创建城池
	for i := x - 5; i <= x+5; i++ {
		for j := y - 5; j <= y+5; j++ {
			posId := global.ToPosition(i, j)
			_, ok := r.posCity[posId]
			if ok {
				return false
			}
		}
	}
	return true
}

func (r *roleCityService) GetCitys(rid int) ([]model.MapRoleCity, error) {
	mrs := make([]data.MapRoleCity, 0)
	mr := &data.MapRoleCity{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("城池查询出错", err)
		return nil, common.New(constant.DBError, "城池查询出错")
	}
	modelMrs := make([]model.MapRoleCity, len(mrs))
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.MapRoleCity))
	}
	return modelMrs, nil
}

func (r *roleCityService) ScanBlock(req *model.ScanBlockReq) []model.MapRoleCity {
	x := req.X
	y := req.Y
	length := req.Length
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return nil
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	rb := make([]model.MapRoleCity, 0)
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			v, ok := r.posCity[posId]
			if ok {
				rb = append(rb, v.ToModel().(model.MapRoleCity))
			}
		}
	}
	return rb

}

func (r *roleCityService) Load() {
	dbCity := make(map[int]*data.MapRoleCity)
	err := db.Engine.Find(dbCity)
	if err != nil {
		log.Println("RoleCityService load role_city table error")
		return
	}

	//转成posCity、roleCity
	for _, v := range dbCity {
		posId := global.ToPosition(v.X, v.Y)
		r.posCity[posId] = v
		_, ok := r.roleCity[v.RId]
		if ok == false {
			r.roleCity[v.RId] = make([]*data.MapRoleCity, 0)
		}
		r.roleCity[v.RId] = append(r.roleCity[v.RId], v)
	}
	//耐久度计算 后续做

}
