package psqlstore

import (
	"log"
	"strconv"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PsqlAccountStore struct {
	*PsqlStore
}

func (s *PsqlStore) Account() store.AccountStore {
	return &PsqlAccountStore{s}
}

func (s *PsqlAccountStore) Get(username string) (model.Account, error) {
	var account model.Account
	rows, err := s.db.Query(`
		SELECT
			id,
			username,
			email,
			win,
			loss,
			disputes,
			steam_id,
			finder_post_ids,
			create_date
		FROM
			account
		WHERE
			username = $1
		LIMIT 1
		;`,
		username,
	)

	if err != nil {
		log.Println("e0010: Failed to find 'account' '" + username + "'")
		log.Println(err)
		return model.Account{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var finderPostIds pq.Int64Array
		err = rows.Scan(
			&account.Id,
			&account.Username,
			&account.Email,
			&account.Win,
			&account.Loss,
			&account.Disputes,
			&account.SteamId,
			&finderPostIds,
			&account.CreateDate,
		)
		if err != nil {
			log.Println("e0011: Failed to populate Account struct'")
			log.Println(err)
			return model.Account{}, err
		}

		account.FinderPostIds = []int{}
		for _, id := range finderPostIds {
			account.FinderPostIds = append(
				account.FinderPostIds,
				int(id),
			)
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
			steam_id,
			finder_post_ids,
			create_date
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
		var finderPostIds pq.Int64Array
		err = rows.Scan(
			&account.Id,
			&account.Username,
			&account.Win,
			&account.Loss,
			&account.Disputes,
			&account.SteamId,
			&finderPostIds,
			&account.CreateDate,
		)
		if err != nil {
			log.Println("e0032: Failed to populate Account struct'")
			log.Println(err)
			return []model.Account{}, err
		}

		account.FinderPostIds = []int{}
		for _, id := range finderPostIds {
			account.FinderPostIds = append(
				account.FinderPostIds,
				int(id),
			)
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PsqlAccountStore) Create(account *model.Account) error {
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
				steam_id = $7,
				finder_post_ids = $8
				
			WHERE
				id = $9
			;`,
			account.Username,
			hashedPass,
			account.Email,
			account.Win,
			account.Loss,
			account.Disputes,
			account.SteamId,
			pq.Array(account.FinderPostIds),
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
				steam_id = $6,
				finder_post_ids = $7
			WHERE
				id = $8
			;`,
			account.Username,
			account.Email,
			account.Win,
			account.Loss,
			account.Disputes,
			account.SteamId,
			pq.Array(account.FinderPostIds),
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

	return nil
}

func (s *PsqlAccountStore) IncrementWin(id int) error {
	_, err := s.db.Exec(`
		UPDATE
			account
		SET
			win = win + 1
		WHERE
			id = $1
	`,
		id,
	)
	if err != nil {
		log.Println("e0046: Failed to increase wins on account")
		log.Println(err)
		return err
	}
	return nil
}

func (s *PsqlAccountStore) IncrementLoss(id int) error {
	_, err := s.db.Exec(`
		UPDATE
			account
		SET
			loss = loss + 1
		WHERE
			id = $1
	`,
		id,
	)
	if err != nil {
		log.Println("e0046: Failed to increase losses on account")
		log.Println(err)
		return err
	}
	return nil
}

func (s *PsqlAccountStore) IncrementDispute(id int) error {
	_, err := s.db.Exec(`
		UPDATE
			account
		SET
			disputes = disputes + 1
		WHERE
			id = $1
	`,
		id,
	)
	if err != nil {
		log.Println("e0046: Failed to increase disputes on account")
		log.Println(err)
		return err
	}
	return nil
}
