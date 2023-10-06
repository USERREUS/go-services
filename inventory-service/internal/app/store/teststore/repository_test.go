package teststore_test

import (
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
	"http-rest-api/internal/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	s := teststore.New()
	temp := model.TestModel(t)
	assert.NoError(t, s.Repository().Create(temp))
	assert.NotNil(t, temp.ID)
}

func TestRepository_FindOne(t *testing.T) {
	s := teststore.New()
	temp := model.TestModel(t)
	s.Repository().Create(temp)
	res, err := s.Repository().FindOne(temp.ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestRepository_FindAll(t *testing.T) {
	s := teststore.New()
	temp := model.TestModel(t)
	_, err := s.Repository().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Repository().Create(temp)
	res, err := s.Repository().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
