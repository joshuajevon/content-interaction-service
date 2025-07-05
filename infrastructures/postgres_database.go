package infrastructures

import (
	"bootcamp-content-interaction-service/config"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(
			fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s",
				conf.Db.Host,
				conf.Db.User,
				conf.Db.Password,
				conf.Db.DbName,
				conf.Db.Port,
				conf.Db.SslMode,
				conf.Db.TimeZone,
			),
		),
			&gorm.Config{},
		)

		if err != nil {
			panic(err)
		}

		dbInstance = &PostgresDatabase{
			db: db,
		}
	})

	return dbInstance
}

func (p *PostgresDatabase) GetInstance() *gorm.DB {
	return dbInstance.db
}
