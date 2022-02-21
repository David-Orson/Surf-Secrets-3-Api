package psqlstore

type Migration struct {
	Up   string
	Down string
}

var migrations = []Migration{
	Migration{
		`CREATE TABLE IF NOT EXISTS account	 (
			id serial,
			username varchar(30) NOT NULL DEFAULT '',
			email varchar(50) NOT NULL,
			password varchar(64) NOT NULL,
			win int NOT NULL DEFAULT 0,
			loss int NOT NULL DEFAULT 0,
			disputes int NOT NULL DEFAULT 0,
			steam_id varchar(20),
			finder_post_ids int[] NOT NULL DEFAULT '{}',
			create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			modify_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);`,
		"DROP TABLE IF EXISTS account;",
	},
	Migration{
		`CREATE TABLE IF NOT EXISTS token (
			id serial,
			account_id int NOT NULL,
			username varchar(30) NOT NULL DEFAULT '',
			token char(100) NOT NULL,
			create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			modify_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY (account_id) REFERENCES account(id)
		);`,
		"DROP TABLE IF EXISTS token;",
	},
	Migration{
		`CREATE TABLE IF NOT EXISTS map (
			id serial,
			name varchar(30) NOT NULL,
			tier int NOT NULL,
			create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			modify_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);`,
		"DROP TABLE IF EXISTS map;",
	},
	Migration{
		`CREATE TABLE IF NOT EXISTS finder_post (
			id serial,
			team int[] NOT NULL,
			time TIMESTAMP NOT NULL,
			maps jsonb NOT NULL,
			is_accepted boolean NOT NULL,
			create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			modify_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);`,
		"DROP TABLE IF EXISTS post;",
	},
	Migration{
		`CREATE TABLE IF NOT EXISTS match (
			id serial,
			team_0 int[] NOT NULL,
			team_1 int[] NOT NULL,
			team_size int NOT NULL,
			time TIMESTAMP NOT NULL,
			maps jsonb NOT NULL,
			result_0 int[] NOT NULL,
			result_1 int[] NOT NULL,
			is_disputed boolean NOT NULL,
			result int NOT NULL,
			create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			modify_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);`,
		"DROP TABLE IF EXISTS match;",
	},
}

func (s *PsqlStore) Up() {
	for _, migration := range migrations {
		s.Exec(migration.Up)
	}
}

func (s *PsqlStore) Down() {
	for i, j := 0, len(migrations)-1; i < j; i, j = i+1, j-1 {
		migrations[i], migrations[j] = migrations[j], migrations[i]
	}
	for _, migration := range migrations {
		s.Exec(migration.Down)
	}
}
