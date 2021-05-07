package log

import (
	"log"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"mitaitech.com/oa/pkg/common/config"
)

var sugarLogger *zap.SugaredLogger

func InitLogger() {
	cfg := config.GetConfig()
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	level := getLevel()

	consoleout := zapcore.Lock(os.Stdout)

	var allCore []zapcore.Core
	if cfg.LogStdout {
		allCore = append(allCore, zapcore.NewCore(encoder, consoleout, level))
	}
	allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, level))

	core := zapcore.NewTee(allCore...)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getLevel() zapcore.Level {
	configLevel := strings.ToLower(config.GetConfig().LogLevel)

	levelMap := map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
	if level, ok := levelMap[configLevel]; ok {
		return level
	}

	log.Fatalf("log.level is invalid. expect one of debug, info, wain, error, dpanic, panic, fatal. but got %s.", configLevel)
	return zapcore.InfoLevel
}

func getEncoder() zapcore.Encoder {
	encoder := config.GetConfig().LogEncoder

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if encoder == "jsonEncoder" {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func getLogWriter() zapcore.WriteSyncer {
	cfg := config.GetConfig()

	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.LogFilename,
		MaxSize:    cfg.LogMaxSize,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAge,
		LocalTime:  cfg.LogLocalTime,
		Compress:   cfg.LogCompress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	sugarLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	sugarLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	sugarLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	sugarLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	sugarLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	sugarLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	sugarLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	sugarLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	sugarLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugarLogger.Fatalf(template, args...)
}
