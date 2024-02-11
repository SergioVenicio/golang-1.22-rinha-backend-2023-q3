package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ja7ad/forker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/config"
	"golang.org/x/net/http2"
)

type Server struct {
	HttpServer forker.Forker
	Config     *config.Config
	DB         *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool, cfg *config.Config) *Server {
	httpServer := http.Server{}
	http2.ConfigureServer(&httpServer, &http2.Server{})
	f := forker.New(&httpServer)
	server := &Server{
		HttpServer: f,
		DB:         db,
		Config:     cfg,
	}

	return server
}

func (s *Server) Server() {
	fmt.Printf("starting with config: %v\n", s.Config)
	addr := "0.0.0.0:" + s.Config.Server.Port
	log.Fatal(s.HttpServer.ListenAndServe(addr))
}
