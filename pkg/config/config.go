package config

import (
	"log"
	"os"

	"github.com/go-chassis/go-archaius"
)

type Config struct {
	IP   string
	Port int

	LogStdout     bool
	LogEncoder    string
	LogLevel      string
	LogFilename   string
	LogMaxSize    int
	LogMaxAge     int
	LogMaxBackups int
	LogLocalTime  bool
	LogCompress   bool
}

var config Config
var configFilename string = "./conf.yaml"

func SetConfigFilename(filename string) {
	configFilename = filename
}

func GetConfig() *Config {
	return &config
}

func InitConfig() {
	err := archaius.Init(archaius.WithRequiredFiles([]string{configFilename}))
	if err != nil {
		log.Fatalln("Init config error:" + err.Error())
	}

	config.IP = archaius.GetString("server.ip", "localhost")
	config.Port = archaius.GetInt("server.port", 8080)

	config.LogStdout = archaius.GetBool("log.stdout", true)
	config.LogEncoder = archaius.GetString("log.encoder", "consoleEncoder") // or jsonEncoder
	config.LogLevel = archaius.GetString("log.level", "debug")              // debug info warn error dpanic panic fatal
	config.LogFilename = archaius.GetString("log.filename", os.Args[0]+".log")
	config.LogMaxSize = archaius.GetInt("log.maxsize", 100)
	config.LogMaxAge = archaius.GetInt("log.maxage", 0)
	config.LogMaxBackups = archaius.GetInt("log.maxbackups", 0)
	config.LogLocalTime = archaius.GetBool("log.localtime", false)
	config.LogCompress = archaius.GetBool("log.compress", true)
}
