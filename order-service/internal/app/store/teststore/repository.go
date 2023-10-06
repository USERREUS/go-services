package teststore

import (
	"order-service/internal/app/model"
	"order-service/internal/app/store"

	"github.com/google/uuid"
)

type Repository struct {
	store   *Store
	records map[string]*model.Model
}

func (r *Repository) Create(m *model.Model) error {
	ID := uuid.New().String()
	r.records[ID] = m

	return nil
}

func (r *Repository) FindOne(id string) (*model.Model, error) {
	m, ok := r.records[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return m, nil
}

func (r *Repository) FindAll() (map[string]*model.Model, error) {
	if len(r.records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return r.records, nil
}
