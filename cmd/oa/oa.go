package main

import (
	"mitaitech.com/oa/pkg/config"
	"mitaitech.com/oa/pkg/log"
)

func main() {
	config.InitConfig()
	log.InitLogger()

	log.Debug("debug message")
	log.Infof("info message, v: %s", "hello log")

}
