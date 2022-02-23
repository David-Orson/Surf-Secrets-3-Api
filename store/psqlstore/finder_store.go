package psqlstore

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
	"github.com/lib/pq"
)

type PsqlFinderStore struct {
	*PsqlStore
}

func (s *PsqlStore) Finder() store.FinderStore {
	return &PsqlFinderStore{s}
}

func (s *PsqlFinderStore) GetPost(id int) (model.FinderPost, error) {
	var finderPost model.FinderPost
	var maps []uint8

	rows, err := s.db.Query(`
		SELECT	
			id,
			team,
			time,
			maps
		FROM
			finder_post
		WHERE
			id = $1
		AND
			is_accepted = false
		LIMIT 1
		;`,
		id,
	)

	if err != nil {
		log.Println("e0040: Failed to find finder post with id '" + strconv.Itoa(id) + "'")
		log.Println(err)
		return model.FinderPost{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var teamIds pq.Int64Array
		err = rows.Scan(
			&finderPost.Id,
			&teamIds,
			&finderPost.Time,
			&maps,
		)
		if err != nil {
			log.Println("e0041: Failed to populate FinderPost struct'")
			log.Println(err)
			return model.FinderPost{}, err
		}
		err = json.Unmarshal([]byte(maps), &finderPost.Maps)
		if err != nil {
			log.Println("Error: Failed to read MeasureRecord features")
			log.Println(err)
			return model.FinderPost{}, err
		}
		finderPost.Team = []int{}
		for _, id := range teamIds {
			finderPost.Team = append(
				finderPost.Team,
				int(id),
			)
		}
	}

	return finderPost, nil
}

func (s *PsqlFinderStore) GetAllPosts() ([]model.FinderPost, error) {
	var finderPosts []model.FinderPost
	var maps []uint8

	rows, err := s.db.Query(`
		SELECT
			id,
			team,
			time,
			maps
		FROM
			finder_post
		WHERE
			is_accepted = false
		AND
			time > NOW()
		;`,
	)

	if err != nil {
		log.Println("e0037: Failed to get all match finder posts")
		log.Println(err)
		return []model.FinderPost{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var finderPost model.FinderPost
		var teamIds pq.Int64Array
		err = rows.Scan(
			&finderPost.Id,
			&teamIds,
			&finderPost.Time,
			&maps,
		)
		if err != nil {
			log.Println("e0038: Failed to populate FinderPost struct'")
			log.Println(err)
			return []model.FinderPost{}, err
		}
		err = json.Unmarshal([]byte(maps), &finderPost.Maps)
		if err != nil {
			log.Println("Error: Failed to read MeasureRecord features")
			log.Println(err)
			return []model.FinderPost{}, err
		}

		finderPost.Team = []int{}
		for _, id := range teamIds {
			finderPost.Team = append(
				finderPost.Team,
				int(id),
			)
		}

		finderPosts = append(finderPosts, finderPost)
	}

	return finderPosts, nil
}

func (s *PsqlFinderStore) CreatePost(finderPost *model.FinderPost) error {
	var id int

	maps, err := json.Marshal(finderPost.Maps)

	err = s.db.QueryRow(`
			INSERT INTO finder_post (
				team,
				time,
				maps,
				is_accepted
			) VALUES (
				$1,
				$2,
				$3,
				false
			)
			RETURNING id
			;`,
		pq.Array(finderPost.Team),
		finderPost.Time,
		maps,
	).Scan(&id)
	if err != nil {
		log.Println("e0039: Failed to create 'finder' row")
		log.Println(err)
		return err
	}

	finderPost.Id = id

	_, err = s.db.Exec(`
		UPDATE
			account
		SET
			finder_post_ids = ARRAY_APPEND(finder_post_ids, $1)
		WHERE
			id = $2
	`,
		finderPost.Id,
		finderPost.Team[0],
	)
	if err != nil {
		log.Println("e0045: Failed to add finder port id to account")
		log.Println(err)
		return err
	}

	return nil
}

func (s *PsqlFinderStore) SetAccepted(id int) error {
	_, err := s.db.Exec(`
			UPDATE
				finder_post
			SET
				is_accepted = true
			WHERE
				id = $1
			;`,
		id,
	)
	if err != nil {
		log.Println("e0044: Failed to set finder post as accepted")
		log.Println(err)
		return err
	}

	return nil
}
