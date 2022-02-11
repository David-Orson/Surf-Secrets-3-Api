package psqlstore

import (
	"log"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

type PsqlFinderStore struct {
	*PsqlStore
}

func (s *PsqlStore) Finder() store.FinderStore {
	return &PsqlFinderStore{s}
}

func (s *PsqlFinderStore) GetAllPosts() ([]model.FinderPost, error) {
	var finders []model.FinderPost
	rows, err := s.db.Query(`
		SELECT
			id,
			team,
			team_size,
			time,
			maps
		FROM
			finder
		;`,
	)

	if err != nil {
		log.Println("e0037: Failed to get all match finder posts")
		log.Println(err)
		return []model.FinderPost{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var finder model.FinderPost
		err = rows.Scan(
			&finder.Id,
			&finder.Team,
			&finder.TeamSize,
			&finder.Time,
			&finder.Maps,
		)
		if err != nil {
			log.Println("e0038: Failed to populate Finder struct'")
			log.Println(err)
			return []model.FinderPost{}, err
		}
		finders = append(finders, finder)
	}

	return finders, nil
}

func (s *PsqlFinderStore) CreatePost(finder *model.FinderPost) error {
	var id int
	err := s.db.QueryRow(`
			INSERT INTO finder (
				team,
				team_size,
				time,
				maps
			) VALUES (
				$1,
				$2,
				$3,
				$4
			)
			RETURNING id
			;`,
		finder.Team,
		finder.TeamSize,
		finder.Time,
		finder.Maps,
	).Scan(&id)
	if err != nil {
		log.Println("e0039: Failed to create 'finder' row")
		log.Println(err)
		return err
	}

	finder.Id = id

	return nil
}
