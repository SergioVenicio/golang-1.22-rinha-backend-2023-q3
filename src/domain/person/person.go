package person

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"apelido"`
	Name      string    `json:"nome"`
	Birthdate time.Time `json:"nascimento"`
	Stack     []string  `json:"stack"`
}

func (p *Person) StackStr() string {
	return strings.Join(p.Stack, ",")
}

func (p *Person) SearchStr() string {
	return p.Nickname + " " + p.Name + " " + p.StackStr()
}

func BuildPerson(
	id string, nick string, name string, birthdate time.Time, stack []string,
) *Person {
	return &Person{
		ID:        id,
		Nickname:  nick,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}

func NewPerson(
	nick string, name string, birthdate time.Time, stack []string,
) *Person {
	return &Person{
		ID:        uuid.NewString(),
		Nickname:  nick,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}
