// Package main is the entry point for the kenketsu-plus API server.
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/akaitigo/kenketsu-plus/api/internal/handler"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var (
		centerRepo    repository.CenterRepo
		donationRepo  repository.DonationRepo
		inventoryRepo repository.InventoryRepo
		subRepo       repository.SubscriptionRepo
	)

	storage := os.Getenv("STORAGE")
	if storage == "postgres" {
		db := mustOpenDB()
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Failed to close database: %v", err)
			}
		}()

		centerRepo = repository.NewPgCenterRepository(db)
		donationRepo = repository.NewPgDonationRepository(db)
		inventoryRepo = repository.NewPgInventoryRepository(db)
		subRepo = repository.NewPgSubscriptionRepository(db)
		log.Println("Using PostgreSQL storage")
	} else {
		centerRepo = repository.NewCenterRepository()
		donationRepo = repository.NewDonationRepository()
		inventoryRepo = repository.NewInventoryRepository()
		subRepo = repository.NewSubscriptionRepository()
		log.Println("Using in-memory storage (set STORAGE=postgres for PostgreSQL)")
	}

	calculator := service.NewDonationCalculator()
	router := handler.NewRouter(centerRepo, donationRepo, inventoryRepo, subRepo, calculator)

	// H-1: full timeout configuration
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Starting server on :%s", port) //nolint:gosec // port is not user-controlled taint
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err) //nolint:gocritic // intentional exit on server failure
	}
}

func mustOpenDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required when STORAGE=postgres") //nolint:gocritic // intentional exit on config error
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := pingDB(db); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}

func pingDB(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
