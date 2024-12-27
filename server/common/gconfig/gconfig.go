package gconfig

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"overall/common/utils/util"
	"path"
)

var (
	GConfig *GlobalConfig
)

func checkConfig() {
	if GConfig.BasePath == "" || GConfig.BasePath == "." || GConfig.BasePath == "./" {
		GConfig.BasePath = util.GetAppWd()
	}
	GConfig.LogConf.CliConf.ClientSocket = path.Join(GConfig.BasePath, GConfig.LogConf.CliConf.ClientSocket)
	GConfig.LogConf.CliConf.ServerSocket = path.Join(GConfig.BasePath, GConfig.LogConf.CliConf.ServerSocket)
}

func GetBashPath() string {
	return GConfig.BasePath
}

func LoadGlobalConfig() error {
	// 1. 决策配置文件路径
	// 1.1 配置文件路径
	//     /etc/global.yaml
	//     ./etc/global.yaml
	//     .global.yaml
	// 2. 加载配置文件

	var configFile *string

	configFile = flag.String("release", "/etc/global.yaml", "The Release global config file")
	err := conf.Load(*configFile, &GConfig)
	if err == nil {
		checkConfig()

		return nil
	}

	logx.Infof("Load release config file FAILED %v\n", err)

	configFile = flag.String("debug", path.Join(util.GetAppWd(), "etc", "global.yaml"), "The Debug global config file")
	err = conf.Load(*configFile, &GConfig)
	if err == nil {
		checkConfig()

		return nil
	}

	logx.Infof("Load debug config file FAILED %v\n", err)

	configFile = flag.String("private", "./global.yaml", "The Private global config file")
	err = conf.Load(*configFile, &GConfig)
	if err == nil {
		checkConfig()

		return nil
	}

	logx.Infof("Load private config file FAILED %v\n", err)

	return err
}
