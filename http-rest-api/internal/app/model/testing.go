package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestProduct(t *testing.T) *Product {
	t.Helper()

	return &Product{
		ID:     uuid.New().String(),
		Name:   "testName",
		Weight: 1000,
	}
}
