package logic

import (
	"encoding/json"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/model/data"
)

var CityFacilityService = &cityFacilityService{}

type cityFacilityService struct {
}

func (c *cityFacilityService) TryCreate(cid, rid int) error {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(cf).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println(err)
		return common.New(constant.DBError, "数据库错误")
	}
	if ok {
		return nil
	}
	list := gameConfig.FacilityConf.List
	fs := make([]data.Facility, len(list))
	for index, v := range list {
		fac := data.Facility{
			Name:         v.Name,
			PrivateLevel: 0,
			Type:         v.Type,
			UpTime:       0,
		}
		fs[index] = fac
	}
	fsJson, err := json.Marshal(fs)
	if err != nil {
		return common.New(constant.DBError, "转json出错")
	}
	cf.RId = rid
	cf.CityId = cid
	cf.Facilities = string(fsJson)
	db.Engine.Table(cf).Insert(cf)
	return nil
}
