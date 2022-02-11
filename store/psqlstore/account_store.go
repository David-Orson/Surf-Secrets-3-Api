package psqlstore

import (
	"errors"
	"log"
	"strconv"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
	"golang.org/x/crypto/bcrypt"
)

type PsqlAccountStore struct {
	*PsqlStore
}

func (s *PsqlStore) Account() store.AccountStore {
	return &PsqlAccountStore{s}
}

func (s *PsqlAccountStore) Get(id int) (model.Account, error) {
	var account model.Account
	rows, err := s.db.Query(`
		SELECT
			id,
			username,
			email,
			win,
			loss,
			disputes,
			steam_id
		FROM
			account
		WHERE
			id = $1
		LIMIT 1
		;`,
		id,
	)

	if err != nil {
		log.Println("e0010: Failed to find 'account' with id '" + strconv.Itoa(id) + "'")
		log.Println(err)
		return model.Account{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&account.Id,
			&account.Username,
			&account.Email,
			&account.Win,
			&account.Loss,
			&account.Disputes,
			&account.SteamId,
		)
		if err != nil {
			log.Println("e0011: Failed to populate Account struct'")
			log.Println(err)
			return model.Account{}, err
		}
	}

	return account, nil
}

func (s *PsqlAccountStore) GetAll() ([]model.Account, error) {
	var accounts []model.Account
	rows, err := s.db.Query(`
		SELECT
			id,
			username,
			
			win,
			loss,
			disputes,
			steam_id
		FROM
			account
		;`,
	)

	if err != nil {
		log.Println("e0025: Failed to get all users")
		log.Println(err)
		return []model.Account{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var account model.Account
		err = rows.Scan(
			&account.Id,
			&account.Username,
			&account.Win,
			&account.Loss,
			&account.Disputes,
			&account.SteamId,
		)
		if err != nil {
			log.Println("e0032: Failed to populate Account struct'")
			log.Println(err)
			return []model.Account{}, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PsqlAccountStore) Create(account *model.Account) error {
	if !s.CheckExists("account", "email", account.Email) {
		hashedPass, _ := bcrypt.GenerateFromPassword(
			[]byte(account.Password),
			10,
		)
		var id int
		err := s.db.QueryRow(`
			INSERT INTO account (
				username,
				email,
				password,
				win,
				loss,
				disputes,
				steam_id
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7
			)
			RETURNING id
			;`,
			account.Username,
			account.Email,
			hashedPass,
			account.Win,
			account.Loss,
			account.Disputes,
			account.SteamId,
		).Scan(&id)
		if err != nil {
			log.Println("e0012: Failed to create 'account' row")
			log.Println(err)
			return err
		}

		account.Id = id
	} else {
		err := errors.New("'account' with email '" + account.Email + "' already exists")
		log.Println(err)
		return err
	}
	return nil
}

func (s *PsqlAccountStore) Update(account *model.Account) error {
	if account.Password != "" {
		hashedPass, _ := bcrypt.GenerateFromPassword(
			[]byte(account.Password),
			10,
		)
		_, err := s.db.Exec(`
			UPDATE
				account
			SET
				username = $1,
				password = $2,
				email = $3,
				win = $4,
				loss = $5,
				disputes = $6,
				steam_id = $7
			WHERE
				id = $8
			;`,
			account.Username,
			hashedPass,
			account.Email,
			account.Win,
			account.Loss,
			account.Disputes,
			account.SteamId,
			account.Id,
		)
		if err != nil {
			log.Println("e0013: Failed to update 'account' row")
			log.Println(err)
			return err
		}
	} else {
		_, err := s.db.Exec(`
			UPDATE
				account
			SET
				username = $1,
				email = $2,
				win = $3,
				loss = $4,
				disputes = $5,
				steam_id = $6
			WHERE
				id = $7
			;`,
			account.Username,
			account.Email,
			account.Win,
			account.Loss,
			account.Disputes,
			account.SteamId,
			account.Id,
		)
		if err != nil {
			log.Println("e0014: Failed to update 'account' row")
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *PsqlAccountStore) Delete(id int) error {
	if s.CheckExists("account", "id", id) {
		s.Token().DeleteAllByAccountId(id)

		_, err := s.db.Exec(`
				DELETE FROM
					account
				WHERE
					id = $1
				;`,
			id,
		)
		if err != nil {
			log.Println("e0015: Failed to delete 'account' with id '" + strconv.Itoa(id) + "'")
			log.Println(err)
			return err
		}
	}
	return nil
}
