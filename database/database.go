package database

import (
	"context"
	"cqs-kanban/config"
	"cqs-kanban/internal/models"
	"cqs-kanban/logger"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Database interface {
	DB() *gorm.DB
	ERPDB() *gorm.DB
	Close() error
	Ping() error
}

type database struct {
	db    *gorm.DB
	erpDB *gorm.DB
}

func NewDatabase(cfg *config.Config, log *logger.AppLogger) (Database, error) {

	logger := NewLogger(log.Logger).LogMode(gormlogger.Info)

	// ERP Database
	erpDB, err := newERPDatabase(cfg.GetERPDatabaseDSN(), logger)
	if err != nil {
		panic(fmt.Sprintf("connect to erp database with err: [%v]", err))
	}

	// Main Database
	db, err := newDatabase(cfg.GetDSN(), logger)
	if err != nil {
		panic(fmt.Sprintf("connect to main database with err: [%v]", err))
	}

	return &database{
		db:    db,
		erpDB: erpDB,
	}, nil
}

func MustNewDatabase(cfg *config.Config, logger *logger.AppLogger) Database {
	db, err := NewDatabase(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return db
}

func newDatabase(dns string, logger gormlogger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger,
	})
	db.AutoMigrate(&models.Report{}, &models.ReportColumn{}, &models.Department{}, &models.Menu{},
		&models.User{}, &models.Operation{}, &models.AccessLog{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	// sqlDB.SetMaxIdleConns(cfg.Database.MaxConnections)
	// sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func (d *database) DB() *gorm.DB {
	return d.db
}

func (d *database) ERPDB() *gorm.DB {
	return d.erpDB
}

func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
