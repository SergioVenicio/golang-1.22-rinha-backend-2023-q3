package database

const QueryPeopleCount = "SELECT COUNT(1) FROM people;"

const QueryInsertPerson = "INSERT INTO people(id, nickname, name, birthdate, stack, search) VALUES ($1, $2, $3, $4, $5, $6);"

const QueryNicknameExists = "SELECT COUNT(1) FROM people WHERE nickname = $1;"

const QueryGetPersonById = "SELECT (id, nickname, name, birthdate, stack) FROM people WHERE id = $1;"

const QuerySearchPersonByTerm = "SELECT id, nickname, name, birthdate, stack FROM people WHERE search LIKE '%' || $1 || '%';"
