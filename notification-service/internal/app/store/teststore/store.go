package teststore

import (
	"notification-service/internal/app/model"
	"notification-service/internal/app/store"
)

type Store struct {
	repository *Repository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Repository() store.Repository {
	if s.repository != nil {
		return s.repository
	}

	s.repository = &Repository{
		store:   s,
		records: make(map[string]*model.Model),
	}

	return s.repository
}
