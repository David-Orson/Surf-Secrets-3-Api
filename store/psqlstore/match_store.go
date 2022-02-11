package psqlstore

import (
	"log"
	"strconv"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

type PsqlMatchStore struct {
	*PsqlStore
}

func (s *PsqlStore) Match() store.MatchStore {
	return &PsqlMatchStore{s}
}

func (s *PsqlMatchStore) Get(id int) (model.Match, error) {
	var match model.Match
	rows, err := s.db.Query(`
		SELECT
			id,
			team_0,
			team_1,
			team_size,
			time,
			maps,
			result_0,
			result_1,
			is_disputed,
			result
		FROM
			match
		WHERE
			id = $1
		LIMIT 1
		;`,
		id,
	)

	if err != nil {
		log.Println("e0029: Failed to find 'match' with id '" + strconv.Itoa(id) + "'")
		log.Println(err)
		return model.Match{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&match.Id,
			&match.Team0,
			&match.Team1,
			&match.TeamSize,
			&match.Time,
			&match.Maps,
			&match.Result0,
			&match.Result1,
			&match.IsDisputed,
			&match.Result,
		)
		if err != nil {
			log.Println("e0033: Failed to populate Match struct'")
			log.Println(err)
			return model.Match{}, err
		}
	}

	return match, nil
}

func (s *PsqlMatchStore) GetAll() ([]model.Match, error) {
	var matchs []model.Match
	rows, err := s.db.Query(`
		SELECT
			id,
			team_0,
			team_1,
			team_size,
			time,
			maps,
			result_0,
			result_1,
			is_disputed,
			result
		FROM
			match
		;`,
	)

	if err != nil {
		log.Println("e0026: Failed to get all users")
		log.Println(err)
		return []model.Match{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var match model.Match
		err = rows.Scan(
			&match.Id,
			&match.Team0,
			&match.Team1,
			&match.TeamSize,
			&match.Time,
			&match.Maps,
			&match.Result0,
			&match.Result1,
			&match.IsDisputed,
			&match.Result,
		)
		if err != nil {
			log.Println("e0034: Failed to populate Match struct'")
			log.Println(err)
			return []model.Match{}, err
		}
		matchs = append(matchs, match)
	}

	return matchs, nil
}

func (s *PsqlMatchStore) GetByAccount(accountId int) ([]model.Match, error) {
	var matchs []model.Match
	// query not tested, account is in either team so filter if true of teams
	rows, err := s.db.Query(`
		SELECT
			id,
			team_0,
			team_1,
			team_size,
			time,
			maps,
			result_0,
			result_1,
			is_disputed,
			result
		FROM
			match
		WHERE
			$1 = ANY team_0
		OR
			$1 = ANY team_1
		;`,
		accountId,
	)

	if err != nil {
		log.Println("e0027: Failed to get all users")
		log.Println(err)
		return []model.Match{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var match model.Match
		err = rows.Scan(
			&match.Id,
			&match.Team0,
			&match.Team1,
			&match.TeamSize,
			&match.Time,
			&match.Maps,
			&match.Result0,
			&match.Result1,
			&match.IsDisputed,
			&match.Result,
		)
		if err != nil {
			log.Println("e0035: Failed to populate Match struct'")
			log.Println(err)
			return []model.Match{}, err
		}
		matchs = append(matchs, match)
	}

	return matchs, nil
}

func (s *PsqlMatchStore) GetDisputesByAccount(accountId int) ([]model.Match, error) {
	var matchs []model.Match
	// query not tested, account is in either team so filter if true of teams
	rows, err := s.db.Query(`
		SELECT
			id,
			team_0,
			team_1,
			team_size,
			time,
			maps,
			result_0,
			result_1,
			is_disputed,
			result
		FROM
			match
		WHERE
			$1 = ANY team_0
		OR
			$1 = ANY team_1
		AND
			is_disputed = true
		;`,
		accountId,
	)

	if err != nil {
		log.Println("e0028: Failed to get all users")
		log.Println(err)
		return []model.Match{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var match model.Match
		err = rows.Scan(
			&match.Id,
			&match.Team0,
			&match.Team1,
			&match.TeamSize,
			&match.Time,
			&match.Maps,
			&match.Result0,
			&match.Result1,
			&match.IsDisputed,
			&match.Result,
		)
		if err != nil {
			log.Println("e0036: Failed to populate Match struct'")
			log.Println(err)
			return []model.Match{}, err
		}
		matchs = append(matchs, match)
	}

	return matchs, nil
}

func (s *PsqlMatchStore) Create(match *model.Match) error {
	var id int
	err := s.db.QueryRow(`
			INSERT INTO match (
				team_0,
				team_1,
				team_size,
				time,
				maps,
				result_0,
				result_1,
				is_disputed,
				result
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7
				$8,
				$9
			)
			RETURNING id
			;`,
		match.Team0,
		match.Team1,
		match.TeamSize,
		match.Time,
		match.Maps,
		match.Result0,
		match.Result1,
		match.IsDisputed,
		match.Result,
	).Scan(&id)
	if err != nil {
		log.Println("e0042: Failed to create 'match' row")
		log.Println(err)
		return err
	}

	match.Id = id

	return nil
}

func (s *PsqlMatchStore) Update(match *model.Match) error {
	_, err := s.db.Exec(`
			UPDATE
				match
			SET
				team_0 = $1,
				team_1 = $2,
				team_size = $3,
				time = $4,
				maps = $5,
				result_0 = $6,
				result_1 = $7,
				is_disputed = $8,
				result = $9
			WHERE
				id = $10
			;`,
		match.Team0,
		match.Team1,
		match.TeamSize,
		match.Time,
		match.Maps,
		match.Result0,
		match.Result1,
		match.IsDisputed,
		match.Result,
		match.Id,
	)
	if err != nil {
		log.Println("e0043: Failed to update 'account' row")
		log.Println(err)
		return err
	}

	return nil
}
