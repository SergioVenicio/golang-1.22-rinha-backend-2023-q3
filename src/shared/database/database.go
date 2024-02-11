package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergiovenicio/rinhaGolang/src/config"
)

var (
	db *pgxpool.Pool
)

func NewDataBase(cfg *config.Config) *pgxpool.Pool {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalln("Unable to parse connection url:", err)
	}

	db, err = pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}
	return db
}
