package db

import (
	"mcp_service/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	DB *gorm.DB
}

var postgreSQLInstance *PostgreSQL

func InitPostgreSQL(dsn string) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		logx.Errorw("初始化PostgreSQL失败", logx.Field("错误", err.Error()))

	}
	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour * 24)
	postgreSQLInstance = &PostgreSQL{DB: db}
	db.AutoMigrate(&model.Rsi{})
	logx.Infof("初始化PostgreSQL成功")
}

func GetPostgreSQL() *PostgreSQL {
	return postgreSQLInstance
}

func (p *PostgreSQL) AutoMigrate(models ...interface{}) error {
	return p.DB.AutoMigrate(models...)
}
