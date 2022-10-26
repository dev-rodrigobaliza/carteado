package database

import (
	"fmt"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func createConnection(cfg *config.App, gc *gorm.Config) (*gorm.DB, error) {
	switch cfg.Database.Type {
	case "postgresql":
		dbConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s connect_timeout=5 TimeZone=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Username,
			cfg.Database.Name,
			cfg.Timezone,
		)
		if cfg.Database.Password != "" {
			dbConnection = fmt.Sprintf("%s password=%s", dbConnection, cfg.Database.Password)
		}
		if cfg.Database.SSL {
			dbConnection = fmt.Sprintf("%s sslmode=require", dbConnection)
		} else {
			dbConnection = fmt.Sprintf("%s sslmode=disable", dbConnection)
		}

		return gorm.Open(postgres.Open(dbConnection), gc)

	case "mssql":
		dbConnection := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
		return gorm.Open(sqlserver.Open(dbConnection), gc)

	case "mysql":
		dbConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
		return gorm.Open(mysql.Open(dbConnection), gc)
	default:
		return nil, errors.ErrInvalidDatabaseType
	}
}
