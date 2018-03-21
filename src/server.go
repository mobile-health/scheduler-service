package main

import (
	"flag"

	"github.com/canhlinh/log4go"

	"github.com/mobile-health/scheduler-service/src/api1"
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/services"
	"github.com/mobile-health/scheduler-service/src/stores"
	"github.com/mobile-health/scheduler-service/src/utils"
	"goji.io"
)

var srv *services.Srv

func Init() {
	var configPath string
	flag.StringVar(&configPath, "conf", "", "path of the config file")
	flag.Parse()

	if configPath == "" {
		configPath = "./conf/config.yaml"
	}

	config.Load("./conf/config.yaml")
	utils.Init("./i18n")
}

func StartServer() {
	log4go.Info("Start server, listen at %s", config.GetConfig().Server.ListenAddress)

	srv = services.NewServer(goji.NewMux(), stores.NewStore())
	api1.Init(srv)
}

func main() {
	Init()
	StartServer()
}
