package login

import (
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/login/controller"
)

var Router = &net.Router{}

func Init() {
	//测试数据库，并且初始化数据库
	db.TestDB()
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)

}
