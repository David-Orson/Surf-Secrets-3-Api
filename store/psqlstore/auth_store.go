package psqlstore

import (
	"log"
	"strings"

	"github.com/David-Orson/Surf-Secrets-3-Api/crypto"
	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
	"golang.org/x/crypto/bcrypt"
)

type PsqlAuthStore struct {
	*PsqlStore
}

func (s *PsqlStore) Auth() store.AuthStore {
	return &PsqlAuthStore{s}
}

func (s *PsqlAuthStore) Login(account *model.Account) (model.Token, error) {
	var hashedPass string
	err := s.db.QueryRow(`
		SELECT
			id,
			username,
			password
		FROM
			account
		WHERE
			email = $1
		;`,
		strings.ToLower(account.Email),
	).Scan(
		&account.Id,
		&account.Username,
		&hashedPass,
	)
	if err != nil {
		log.Println("e0016: Failed to find 'account' with matching id and hashed password")
		log.Println(err)
		return model.Token{}, err
	}

	var tokenModel model.Token
	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPass),
		[]byte(account.Password),
	)
	if err != nil {
		log.Println("e0017: Failed login attempt by '" + account.Email + "'")
		log.Println(err)
		return model.Token{}, err
	}

	token, err := crypto.GenerateToken()
	if err != nil {
		log.Println("e0018: Failed to generate a token")
		log.Println(err)
		return model.Token{}, err
	}

	_, err = s.db.Exec(`
		INSERT INTO token (
			token,
			username,
			account_id
		) VALUES (
			$1,
			$2,
			$3
		)
		;`,
		token,
		account.Username,
		account.Id,
	)
	if err != nil {
		log.Println("e0019: Failed to create 'token' row")
		log.Println(err)
		return model.Token{}, err
	}

	tokenModel.AccountId = account.Id
	tokenModel.Username = account.Username
	tokenModel.Token = token

	return tokenModel, nil
}
