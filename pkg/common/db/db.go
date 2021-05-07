package db

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mitaitech.com/oa/pkg/common/config"
	"mitaitech.com/oa/pkg/common/log"
)

const (
	DB_MYSQL    string = "mysql"
	DB_SQLITE   string = "sqlite"
	DB_POSTGRES string = "postgres"
)

var (
	globalDB *gorm.DB

	injectors []func(db *gorm.DB)
)

func InitDB() {
	cfg := config.GetConfig()

	// 连接数据库前初始化Database
	CreateDatabase()

	var ormLogger logger.Interface
	if cfg.DBDebug {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	var db *gorm.DB
	var err error

	var dsn string = cfg.DBDSN
	log.Debugf("db dsn: %s", dsn)
	switch cfg.DBType {
	case DB_MYSQL:
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: ormLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "tb_", // 表名前缀，`User` 对应的表名是 `tb_users`
			},
		})
	case DB_SQLITE:
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: ormLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "tb_", // 表名前缀，`User` 对应的表名是 `tb_users`
			},
		})
	case DB_POSTGRES:
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: ormLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "tb_", // 表名前缀，`User` 对应的表名是 `tb_users`
			},
		})
	default:
		log.Fatalf("db type is invalid, expect one of mysql, sqlite, postgres, but got %s.", cfg.DBType)
	}

	if err != nil {
		log.Fatal(err)
	}

	idb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	idb.SetMaxIdleConns(cfg.DBMaxIdleConns)
	idb.SetMaxOpenConns(cfg.DBMaxOpenConns)

	registerCallback(db)
	callInjector(db)
	globalDB = db

	log.Info("db connected success")
}

func CreateDatabase() {
	cfg := config.GetConfig()
	switch cfg.DBType {
	case DB_MYSQL:
		CreateDatabaseMysql()
	case DB_POSTGRES:
		CreateDatabasePostgres()
	case DB_SQLITE:
		return
	default:
		log.Fatalf("db type is invalid, expect one of mysql, sqlite, postgres, but got %s.", cfg.DBType)
	}
}

func CreateDatabaseMysql() {
	cfg := config.GetConfig()
	slashIndex := strings.LastIndex(cfg.DBDSN, "/")
	dsn := cfg.DBDSN[:slashIndex+1]
	dbNameAndParam := cfg.DBDSN[slashIndex+1:]

	markIndex := strings.Index(dbNameAndParam, "?")
	dbName := dbNameAndParam[:markIndex+1]
	param := dbNameAndParam[markIndex+1:]

	dsn = dsn + "?" + param
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		log.Fatal(err)
	}

	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4;",
		dbName,
	)

	err = db.Exec(createSQL).Error
	if err != nil {
		log.Fatal(err)
	}
}

func CreateDatabasePostgres() {
	// TODO
}

func RegisterInjector(f func(*gorm.DB)) {
	injectors = append(injectors, f)
}

func callInjector(db *gorm.DB) {
	for _, v := range injectors {
		v(db)
	}
}

type ctxTransactionKey struct{}

func CtxWithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxTransactionKey{}, tx)
}

type txImpl struct{}

func NewTxImpl() *txImpl {
	return &txImpl{}
}

func (*txImpl) Transaction(ctx context.Context, fc func(txctx context.Context) error) error {
	db := globalDB.WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txctx := CtxWithTransaction(ctx, tx)
		return fc(txctx)
	})
}

// 如果使用跨模型事务则传参
func GetDB(ctx context.Context) *gorm.DB {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		tx, ok := iface.(*gorm.DB)
		if !ok {
			log.Panicf("unexpect context value type: %s", reflect.TypeOf(tx))
			return nil
		}

		return tx
	}

	return globalDB.WithContext(ctx)
}

// 自动初始化表结构
func SetupTableModel(db *gorm.DB, model interface{}) {
	cfg := config.GetConfig()

	if cfg.DBAutoMigrate {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// https://github.com/ulid/spec
// uuid sortable by time
func NewUlid() string {
	now := time.Now()
	return ulid.MustNew(ulid.Timestamp(now), ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)).String()
}

func registerCallback(db *gorm.DB) {
	// 自动添加uuid
	err := db.Callback().Create().Before("gorm:create").Register("uuid", func(db *gorm.DB) {
		db.Statement.SetColumn("id", NewUlid())
	})
	if err != nil {
		log.Panicf("err: %+v", err)
	}
}

func WithOffsetLimit(db *gorm.DB, offset, limit int) *gorm.DB {
	if offset > 0 {
		db = db.Offset(offset)
	}

	if limit > 0 {
		db = db.Limit(limit)
	}

	return db
}
