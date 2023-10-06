package teststore

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
)

type Store struct {
	productRepository *ProductRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Product() store.ProductRepository {
	if s.productRepository != nil {
		return s.productRepository
	}

	s.productRepository = &ProductRepository{
		store:    s,
		products: make(map[string]*model.Product),
	}

	return s.productRepository
}
