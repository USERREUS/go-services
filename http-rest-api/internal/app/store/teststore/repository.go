package teststore

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"

	"github.com/google/uuid"
)

type ProductRepository struct {
	store    *Store
	products map[string]*model.Product
}

func (r *ProductRepository) Create(p *model.Product) error {
	p.ID = uuid.New().String()
	r.products[p.ID] = p

	return nil
}

func (r *ProductRepository) FindOne(id string) (*model.Product, error) {
	p, ok := r.products[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return p, nil
}

func (r *ProductRepository) FindAll() (map[string]*model.Product, error) {
	if len(r.products) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return r.products, nil
}
