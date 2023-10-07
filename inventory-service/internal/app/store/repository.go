package store

import "inventory/internal/app/model"

type Repository interface {
	Create(*model.Model) error
	Update(*model.Model) error
	FindOne(string) (*model.Model, error)
	FindAll() (map[string]*model.Model, error)
}
