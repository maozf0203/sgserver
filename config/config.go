package config

import (
	"errors"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

// 加载此文件的时候会执行初始化方法
func init() {
	//获取当前程序目录
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentPath + configFile

	//命令行是否制定配置文件目录，如果指定则使用命令行指定的目录
	args := os.Args
	if len(args) > 1 {
		dir := args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}
	if !fileExists(configPath) {
		panic(errors.New("配置文件不存在"))
	}
	//文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错", err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)

}
