package main

import (
	"mitaitech.com/oa/pkg/common/config"
	"mitaitech.com/oa/pkg/common/log"
)

func main() {
	config.InitConfig()
	log.InitLogger()

	log.Debug("debug message")
	log.Infof("info message, v: %s", "hello log")

}
