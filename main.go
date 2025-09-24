package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/stonejianbu/memo-assistant/config"
	"github.com/stonejianbu/memo-assistant/router"
)

var confFilename string

func init() {
	flag.StringVar(&confFilename, "conf", "./config/config.yaml", "config path, eg: -conf config.yaml")
}
func main() {
	flag.Parse()
	// init config
	config.Init(confFilename)
	addr := config.Cfg.Server.Addr
	engine := router.SetupRouter()
	logrus.Infof("start to serve at %s, name: %s", addr, config.Cfg.Server.Name)
	if err := engine.Run(addr); err != nil {
		logrus.Fatalln(err)
	}
}
