package psqlstore

import (
	"log"
	"strconv"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
)

// PsqlTokenStore receives a pointer to an PsqlStore.
type PsqlTokenStore struct {
	*PsqlStore
}

// Returns the a pointer to a PsqlTokenStore.
func (s *PsqlStore) Token() store.TokenStore {
	return &PsqlTokenStore{s}
}

func (s *PsqlTokenStore) DeleteAllByAccountId(accountId int) error {
	_, err := s.db.Exec(`
		DELETE FROM
			token
		WHERE
			account_id=$1
		;`,
		accountId,
	)
	if err != nil {
		log.Println("Error: Failed to delete 'token' rows with account_id '" + strconv.Itoa(accountId) + "'")
		log.Println(err)
		return err
	}

	return nil
}
