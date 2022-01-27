package psql

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/unifyi/creme-brulee/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORM(ctx context.Context, baseConf *config.BaseConfig, psqlConf *config.PsqlConfig) *gorm.DB {
	log := ctxlogrus.Extract(ctx)

	dbLogLevel := logger.Silent
	if baseConf.LogLevel.String() == logrus.DebugLevel.String() {
		dbLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: psqlConf.GetDataSourcePSQL().String(),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
