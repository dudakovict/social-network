package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dudakovict/social-network/business/data/post/dbschema"
	"github.com/dudakovict/social-network/business/sys/database"
)

func main() {
	err := migrate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func seed() error {
	cfg := database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost:5433",
		Name:         "postgres",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	}

	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}

func migrate() error {
	cfg := database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost:5433",
		Name:         "postgres",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	}

	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")

	return seed()
}
