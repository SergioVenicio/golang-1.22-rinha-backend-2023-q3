package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/sergiovenicio/rinhaGolang/src/domain/person"
	"github.com/sergiovenicio/rinhaGolang/src/shared/repositories"
)

var ErrInvalidDto = errors.New("invalid dto")
var schemaValidator = validator.New()

type CreatePersonRequest struct {
	Nickname  string   `json:"apelido" validate:"required,max=32"`
	Name      string   `json:"nome" validate:"required,max=100"`
	Birthdate string   `json:"nascimento" validate:"required,datetime=2006-01-02"`
	Stack     []string `json:"stack" validate:"dive,max=32"`
}

func (c *CreatePersonRequest) Validate() error {
	if c.Nickname == "" || len(c.Nickname) > 32 {
		return ErrInvalidDto
	}

	if len(c.Name) > 100 {
		return ErrInvalidDto
	}

	for _, tech := range c.Stack {
		if len(tech) > 32 {
			return ErrInvalidDto
		}
	}

	return nil
}

func PeopleCreate(repo *repositories.PeopleRepository, w http.ResponseWriter, r *http.Request) {
	var dto CreatePersonRequest
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	err := json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dto.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err := schemaValidator.Struct(&dto); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	birthdate, err := time.Parse("2006-01-02", dto.Birthdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newPerson := person.NewPerson(
		dto.Nickname,
		dto.Name,
		birthdate,
		dto.Stack,
	)
	if err := repo.Create(*newPerson); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(newPerson)
	w.Write(resp)
}

func PeopleCount(repo *repositories.PeopleRepository, w http.ResponseWriter, r *http.Request) {
	total, _ := repo.Count()
	w.WriteHeader(http.StatusAccepted)
	data, _ := json.Marshal(total)
	w.Write(data)
}

func GetPerson(repo *repositories.PeopleRepository, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	person, err := repo.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(person)
	w.Write(data)
}

func SearchPerson(repo *repositories.PeopleRepository, w http.ResponseWriter, r *http.Request) {
	var people []person.Person
	search := strings.ToLower((r.URL.Query().Get("t")))
	if search == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	people, _ = repo.SearchPerson(search)
	data, _ := json.Marshal(people)
	w.Write(data)
}
