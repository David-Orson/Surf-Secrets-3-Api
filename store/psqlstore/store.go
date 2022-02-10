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
		log.Println(connectionString)
		log.Println(err)
		log.Println("0002: unable to use data source name; this will not be a connection error, but a DSN parse error or another initialisation error.")
		return nil, err
	}

	err = TestDatabase(s)
	if err != nil {
		log.Println("0003: Shutting down")
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
		log.Println("0004: Failed to connect to DB")
	}
	return nil
}

func (s *PsqlStore) Exec(query string) {
	if _, err := s.db.Exec(query); err != nil {
		log.Println(err)
		log.Println(query)
		log.Println("0008: The above query failed to be executed")
	}
}
