package dbsetup

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/unifyi/creme-brulee/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGORM(ctx context.Context, baseConf *config.BaseConfig, psqlConf *config.PsqlConfig) *gorm.DB {
	log := ctxlogrus.Extract(ctx)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: psqlConf.GetDataSourcePSQL().String(),
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if baseConf.LogLevel == logrus.DebugLevel.String() {
		db = db.Debug()
	}

	return db
}
