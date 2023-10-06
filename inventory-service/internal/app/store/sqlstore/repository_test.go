package sqlstore_test

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
	"http-rest-api/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("inventory")

	s := sqlstore.New(db)
	p := model.TestModel(t)
	assert.NoError(t, s.Repository().Create(p))
	assert.NotNil(t, p.ID)
}

func TestRepository_FindOne(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("inventory")

	s := sqlstore.New(db)
	p1 := model.TestModel(t)
	s.Repository().Create(p1)
	p2, err := s.Repository().FindOne(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
}

func TestProductRepository_FindAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("inventory")

	s := sqlstore.New(db)
	p := model.TestModel(t)
	_, err := s.Repository().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Repository().Create(p)
	temp, err := s.Repository().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, temp)
}
