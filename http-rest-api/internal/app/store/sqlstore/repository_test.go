package sqlstore_test

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
	"http-rest-api/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("products")

	s := sqlstore.New(db)
	p := model.TestProduct(t)
	assert.NoError(t, s.Product().Create(p))
	assert.NotNil(t, p.ID)
}

func TestProductRepository_FindOne(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("products")

	s := sqlstore.New(db)
	p1 := model.TestProduct(t)
	s.Product().Create(p1)
	p2, err := s.Product().FindOne(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
}

func TestProductRepository_FindAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("products")

	s := sqlstore.New(db)
	p := model.TestProduct(t)
	_, err := s.Product().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Product().Create(p)
	temp, err := s.Product().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, temp)
}
