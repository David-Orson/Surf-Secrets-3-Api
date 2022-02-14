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
		err = rows.Scan(
			&finderPost.Id,
			pq.Array(&finderPost.Team),
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
				maps
			) VALUES (
				$1,
				$2,
				$3
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

	return nil
}
