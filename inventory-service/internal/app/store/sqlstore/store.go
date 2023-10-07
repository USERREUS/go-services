package sqlstore

import (
	"database/sql"
	"inventory/internal/app/store"
)

type Store struct {
	db         *sql.DB
	repository *Repository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Repository() store.Repository {
	if s.repository != nil {
		return s.repository
	}

	s.repository = &Repository{
		store: s,
	}

	return s.repository
}
