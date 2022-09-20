package model

type MyGeneralReq struct {
}

type MyGeneralRsp struct {
	Generals []General `json:"generals"`
}
