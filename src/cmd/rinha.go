package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/config"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/shared/controllers"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/shared/database"
	sHttp "github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/shared/http"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/shared/repositories"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/shared/workers"
)

func main() {
	godotenv.Load(".env")
	uuid.EnableRandPool()

	jobs := make(workers.JobQueue)

	cfg := config.NewConfig()
	db := database.NewDataBase(cfg)
	srv := sHttp.NewServer(db, cfg)
	dispatcher := workers.NewDispatcher(db, cfg, jobs)
	peopleRepository := repositories.NewPersonRepository(db, jobs)

	http.HandleFunc("GET /contagem-pessoas", func(w http.ResponseWriter, r *http.Request) {
		controllers.PeopleCount(peopleRepository, w, r)
	})
	http.HandleFunc("POST /pessoas", func(w http.ResponseWriter, r *http.Request) {
		controllers.PeopleCreate(peopleRepository, w, r)
	})
	http.HandleFunc("GET /pessoas/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetPerson(peopleRepository, w, r)
	})
	http.HandleFunc("GET /pessoas", func(w http.ResponseWriter, r *http.Request) {
		controllers.SearchPerson(peopleRepository, w, r)
	})

	go dispatcher.Run()
	srv.Server()
}
