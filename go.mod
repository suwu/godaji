module mitaitech.com/oa

go 1.15

replace mitaitech.com/oa => ./

require (
	github.com/go-chassis/go-archaius v1.5.3
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/oklog/ulid/v2 v2.0.2
	github.com/win5do/go-lib v0.0.0-20210322065409-edc6813f5414
	go.uber.org/zap v1.16.0
	gorm.io/driver/mysql v1.0.6
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.9
)
