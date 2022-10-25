package database

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/config"
	"github.com/dev-rodrigobaliza/carteado/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// New creates a new database instance
func New(cfg *config.App) (*Database, error) {
	// create database connection
	var log logger.Interface
	if cfg.Debug {
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default.LogMode(logger.Silent)
	}
	db, err := createConnection(cfg, &gorm.Config{
		Logger: log,
	})
	if err != nil {
		return nil, err
	}
	// create database pool
	sqlDB, err := db.DB()
	if err != nil {
		if sqlDB != nil {
			sqlDB.Close()
		}
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.Database.Pool.Idle)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.Database.Pool.Open)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.Pool.Lifetime) * time.Second)

	database := &Database{
		DB: db,
	}

	err = database.migrate()
	if err != nil {
		return nil, err
	}

	err = database.seed()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) migrate() error {
	err := d.DB.AutoMigrate(
		&model.User{},
		&model.Login{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) seed() error {
	var count int64 = 0
	d.DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	user := model.NewUser("admin", "admin@localhost.com", "123", true)
	now := time.Now()
	user.ConfirmedEmailAt = &now

	return d.DB.Save(user).Error
}
