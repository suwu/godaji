package oa

import (
	"log"

	"github.com/go-chassis/go-archaius"
)

type Config struct {
	IP   string
	Port int
}

var config Config

func InitConfig() {
	err := archaius.Init(archaius.WithRequiredFiles([]string{"./conf/conf.yaml"}))
	if err != nil {
		log.Fatel("Error:" + err.Error())
	}

	config.IP = archaius.GetString("server.ip", "localhost")
	config.Port = archaius.GetInt("server.port", "8080")
}
