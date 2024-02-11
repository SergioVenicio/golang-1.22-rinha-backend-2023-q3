package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergiovenicio/rinhaGolang/src/domain/person"
	"github.com/sergiovenicio/rinhaGolang/src/shared/database"
	"github.com/sergiovenicio/rinhaGolang/src/shared/workers"
)

type PeopleRepository struct {
	jobQueue workers.JobQueue
	db       *pgxpool.Pool
}

func (p *PeopleRepository) Count() (int64, error) {
	var total int64

	err := p.db.
		QueryRow(
			context.Background(),
			database.QueryPeopleCount,
		).
		Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (p *PeopleRepository) FindByNickName(nickname string) bool {
	var count int64
	err := p.db.
		QueryRow(
			context.Background(),
			database.QueryNicknameExists,
			nickname,
		).
		Scan(&count)
	if err != nil || count == 0 {
		return false
	}

	return true
}

func (p *PeopleRepository) Create(person person.Person) error {
	if p.FindByNickName(person.Nickname) {
		return errors.New("invalid nickname")
	}

	p.jobQueue <- workers.Job{
		Payload: &person,
	}
	return nil
}

func (p *PeopleRepository) GetById(id string) (person.Person, error) {
	var person person.Person
	ctx := context.Background()
	p.db.QueryRow(
		ctx,
		database.QueryGetPersonById,
		id,
	).Scan(&person)
	if person.ID == "" {
		return person, errors.New("person not found")
	}
	return person, nil
}

func (p *PeopleRepository) SearchPerson(term string) ([]person.Person, error) {
	data := make([]person.Person, 0)
	ctx := context.Background()
	rows, err := p.db.Query(
		ctx,
		database.QuerySearchPersonByTerm,
		term,
	)
	if err != nil {
		return make([]person.Person, 0), errors.New("error on search by term")
	}
	defer rows.Close()

	for rows.Next() {
		var eachPerson person.Person
		var strStack string
		err := rows.Scan(
			&eachPerson.ID,
			&eachPerson.Nickname,
			&eachPerson.Name,
			&eachPerson.Birthdate,
			&strStack,
		)
		if err != nil {
			return make([]person.Person, 0), errors.New("error scanning person")
		}

		eachPerson.Stack = strings.Split(strStack, ",")
		data = append(data, eachPerson)
	}

	return data, nil
}

func NewPersonRepository(db *pgxpool.Pool, queue workers.JobQueue) *PeopleRepository {
	return &PeopleRepository{
		db:       db,
		jobQueue: queue,
	}
}
