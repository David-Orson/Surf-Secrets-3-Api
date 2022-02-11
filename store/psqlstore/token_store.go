package psqlstore

import (
	"log"
	"strconv"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

type PsqlTokenStore struct {
	*PsqlStore
}

func (s *PsqlStore) Token() store.TokenStore {
	return &PsqlTokenStore{s}
}

func (s *PsqlTokenStore) GetAll() ([]model.Token, error) {
	var tokens []model.Token
	rows, err := s.db.Query(`
		SELECT
			account_id,
			token
		FROM
			token
		;`,
	)
	if err != nil {
		log.Println("e0021: Failed to retrieve 'token' rows")
		log.Println(err)
		return []model.Token{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var token model.Token
		err = rows.Scan(
			&token.AccountId,
			&token.Token,
		)
		if err != nil {
			log.Println("e0022: Failed to populate Token struct")
			log.Println(err)
			return []model.Token{}, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
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
		log.Println("e0023: Failed to delete 'token' rows with account_id '" + strconv.Itoa(accountId) + "'")
		log.Println(err)
		return err
	}

	return nil
}
