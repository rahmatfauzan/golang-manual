package database

import (
	"context"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/rahmatfauzan/golang-manual/internal/config"
)

func ConnectDB(cfg *config.Config) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.DB_URL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v\nURL: %s", err, cfg.DB_URL)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.PingContext(ctx); err != nil {
		log.Fatalf("Database tidak merespon ping: %v", err)
	}

	log.Println("✅ Database connected")

	return db
}
