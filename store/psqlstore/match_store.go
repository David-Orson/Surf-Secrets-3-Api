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
			&match.Result,
		)
		if err != nil {
			log.Println("e0011: Failed to populate Match struct'")
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
			&match.Result,
		)
		if err != nil {
			log.Println("e0011: Failed to populate Match struct'")
			log.Println(err)
			return []model.Match{}, err
		}
		matchs = append(matchs, match)
	}

	return matchs, nil
}

func (s *PsqlMatchStore) GetByAccount(AccountId int) ([]model.Match, error) {
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
		AccountId,
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
			&match.Result,
		)
		if err != nil {
			log.Println("e0011: Failed to populate Match struct'")
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
			&match.Result,
		)
		if err != nil {
			log.Println("e0011: Failed to populate Match struct'")
			log.Println(err)
			return []model.Match{}, err
		}
		matchs = append(matchs, match)
	}

	return matchs, nil
}
