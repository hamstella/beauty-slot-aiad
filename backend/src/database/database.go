package database

import (
	"app/src/config"
	"app/src/model"
	"app/src/utils"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbHost, dbName string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, config.DBUser, config.DBPassword, dbName, config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", err)
		// Return nil database connection but don't panic
		return nil
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Failed to get database instance: %+v", errDB)
		return nil
	}

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	// Auto-migrate beauty salon models
	err = db.AutoMigrate(
		&model.Customer{},
		&model.Staff{},
		&model.Menu{},
		&model.Option{},
		&model.Label{},
		&model.Shift{},
		&model.Reservation{},
		&model.ReservationMenu{},
		&model.ReservationOption{},
		&model.AuditLog{},
		&model.NotificationLog{},
	)
	if err != nil {
		utils.Log.Errorf("Failed to auto-migrate models: %+v", err)
	} else {
		utils.Log.Info("Database models auto-migrated successfully")
	}

	return db
}
