package database

import (
	"fmt"
	"log"
	"os"
	"testovoye/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func InitPostgreSQL() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		os.Getenv("PGSQL_HOST"), os.Getenv("PGSQL_PORT"), os.Getenv("PGSQL_USER"), os.Getenv("PGSQL_PASSWORD"), os.Getenv("PGSQL_DB"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&domain.Person{}); err != nil {
		return nil, err
	}

	log.Println("[INFO] [PGSQL] Connected to PostgreSQL")

	return db, nil
}

func CloseDB() {
	if sqlDB, err := db.DB(); err != nil {
		log.Fatalf("[FATAL] [PGSQL] %s", err.Error())
	} else {
		sqlDB.Close()
	}
}
