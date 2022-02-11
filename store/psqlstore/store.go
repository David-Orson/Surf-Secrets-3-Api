package psqlstore

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/David-Orson/Surf-Secrets-3-Api/config"
	"github.com/David-Orson/Surf-Secrets-3-Api/store"
)

var _ store.Store = &PsqlStore{}

type PsqlStore struct {
	db *sql.DB
}

func Open(jsonConfig string) (*PsqlStore, error) {
	configuration := config.LoadConfigDb(jsonConfig)
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configuration.Host,
		configuration.Port,
		configuration.Username,
		configuration.Password,
		configuration.DbName,
	)

	var s PsqlStore
	var err error
	s.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Println(err)
		log.Println("e0002: unable to use data source name; this will not be a connection error, but a DSN parse error or another initialisation error.")
		return nil, err
	}

	err = TestDatabase(s)
	if err != nil {
		log.Println("e0003: Shutting down")
		os.Exit(1)
	}

	s.Up()

	return &s, nil
}

func (s *PsqlStore) Close() {
	s.Close()
}

func TestDatabase(s PsqlStore) error {
	err := s.db.Ping()
	if err != nil {
		log.Println(err)
		log.Println("e0004: Failed to connect to DB")
	}
	return nil
}

func (s *PsqlStore) Exec(query string) {
	if _, err := s.db.Exec(query); err != nil {
		log.Println(err)
		log.Println(query)
		log.Println("e0008: The above query failed to be executed")
	}
}

func (s *PsqlStore) CheckExists(table string, column string, field interface{}) bool {
	rows, err := s.db.Query(
		"SELECT null FROM "+table+" WHERE "+column+"=$1",
		field,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	return count > 0
}
