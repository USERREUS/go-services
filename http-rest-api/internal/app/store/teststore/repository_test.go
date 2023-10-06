package teststore_test

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
	"http-rest-api/internal/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository_Create(t *testing.T) {
	s := teststore.New()
	p := model.TestProduct(t)
	assert.NoError(t, s.Product().Create(p))
	assert.NotNil(t, p.ID)
}

func TestProductRepository_FindOne(t *testing.T) {
	s := teststore.New()
	p1 := model.TestProduct(t)
	s.Product().Create(p1)
	p2, err := s.Product().FindOne(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
}

func TestProductRepository_FindAll(t *testing.T) {
	s := teststore.New()
	p := model.TestProduct(t)
	_, err := s.Product().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Product().Create(p)
	temp, err := s.Product().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, temp)
}
