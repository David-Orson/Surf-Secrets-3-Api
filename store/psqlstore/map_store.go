package psqlstore

import (
	"log"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

type PsqlMapStore struct {
	*PsqlStore
}

func (s *PsqlStore) Map() store.MapStore {
	return &PsqlMapStore{s}
}

func (s *PsqlMapStore) GetAll() ([]model.Map, error) {
	var maps []model.Map
	rows, err := s.db.Query(`
		SELECT
			id,
			name,
			tier
		FROM
			map
		;`,
	)

	if err != nil {
		log.Println("e0030: Failed to get all users")
		log.Println(err)
		return []model.Map{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var surfMap model.Map
		err = rows.Scan(
			&surfMap.Id,
			&surfMap.Name,
			&surfMap.Tier,
		)
		if err != nil {
			log.Println("e0031: Failed to populate Map struct'")
			log.Println(err)
			return []model.Map{}, err
		}
		maps = append(maps, surfMap)
	}

	return maps, nil
}
