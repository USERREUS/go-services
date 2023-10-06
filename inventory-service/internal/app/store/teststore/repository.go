package teststore

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"

	"github.com/google/uuid"
)

type Repository struct {
	store   *Store
	records map[string]*model.Model
}

func (r *Repository) Create(p *model.Model) error {
	p.ID = uuid.New().String()
	r.records[p.ID] = p

	return nil
}

func (r *Repository) FindOne(id string) (*model.Model, error) {
	p, ok := r.records[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return p, nil
}

func (r *Repository) FindAll() (map[string]*model.Model, error) {
	if len(r.records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return r.records, nil
}
