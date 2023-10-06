package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestModel(t *testing.T) *Model {
	t.Helper()

	return &Model{
		ID:    uuid.New().String(),
		Name:  "testName",
		Count: 100,
		Cost:  50000,
	}
}
