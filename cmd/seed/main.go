package main

import (
	"log"

	"supermarket-backend/infrastructure/db"
	"supermarket-backend/internal/config"
)

func main() {
	cfg := config.Load()

	database, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := seedBranches(database); err != nil {
		log.Fatal(err)
	}

	if err := seedPositions(database); err != nil {
		log.Fatal(err)
	}

	if err := seedRoles(database); err != nil {
		log.Fatal(err)
	}

	if err := seedUnits(database); err != nil {
		log.Fatal(err)
	}

	if err := seedBrands(database); err != nil {
		log.Fatal(err)
	}

	if err := seedSKUs(database); err != nil {
		log.Fatal(err)
	}

	if err := seedStocks(database); err != nil {
		log.Fatal(err)
	}

	if err := seedUsers(database, cfg); err != nil {
		log.Fatal(err)
	}

	log.Println("seed completed")
}
