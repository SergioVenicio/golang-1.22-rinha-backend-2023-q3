package config

import (
	"os"
	"strconv"

	"github.com/caarlos0/env/v10"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Server struct {
	Port string
}

type Worker struct {
	Size         int
	BatchMaxSize int
	IntsertTime  int
	MaxWorkers   int
}

type Config struct {
	Database Database
	Server   Server
	Worker   Worker
}

func NewConfig() *Config {
	workerSize, err := strconv.Atoi(os.Getenv("WORKER_SIZE"))
	if err != nil {
		workerSize = 1
	}
	batchMaxSize, err := strconv.Atoi(os.Getenv("WORKER_BATCH_MAX_SIZE"))
	if err != nil {
		batchMaxSize = 1000
	}

	insertTime, err := strconv.Atoi(os.Getenv("WORKER_INSERT_TIME_SECONDS"))
	if err != nil {
		insertTime = 1
	}

	maxWorkers, err := strconv.Atoi(os.Getenv("WORKER_MAX_WORKERS"))
	if err != nil {
		maxWorkers = 1
	}
	cfg := Config{
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PWD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Server: Server{
			Port: os.Getenv("HTTP_PORT"),
		},
		Worker: Worker{
			Size:         workerSize,
			BatchMaxSize: batchMaxSize,
			IntsertTime:  insertTime,
			MaxWorkers:   maxWorkers,
		},
	}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
