package sqlstore

import (
	"database/sql"
	"fmt"
	"inventory/internal/app/model"
	"inventory/internal/app/store"

	"github.com/google/uuid"
)

const (
	table = "inventory"
)

type Repository struct {
	store *Store
}

func (r *Repository) Create(m *model.Model) error {
	m.ID = uuid.New().String()
	return r.store.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s (id, name, count, cost) VALUES ($1, $2, $3, $4) RETURNING id", table),
		m.ID,
		m.Name,
		m.Count,
		m.Cost,
	).Scan(&m.ID)
}

func (r *Repository) Update(m *model.Model) error {
	m.ID = uuid.New().String()
	return r.store.db.QueryRow(
		fmt.Sprintf("UPDATE %s SET name = $2, count = $3, cost = $4 WHERE id = $1 RETURNING id", table), ///FIX??
		m.ID,
		m.Name,
		m.Count,
		m.Cost,
	).Scan(&m.ID)
}

func (r *Repository) FindOne(id string) (*model.Model, error) {
	p := &model.Model{}
	if err := r.store.db.QueryRow(
		fmt.Sprintf("SELECT id, name, count, cost FROM %s WHERE id = $1", table),
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Count,
		&p.Cost,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return p, nil
}

func (r *Repository) FindAll() (map[string]*model.Model, error) {
	records := make(map[string]*model.Model)
	m := &model.Model{}

	rows, err := r.store.db.Query("SELECT id, name, count, cost FROM inventory")
	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Name, &m.Count, &m.Cost); err != nil {
			return nil, store.ErrRecordNotFound
		}

		records[m.ID] = m
	}

	if err := rows.Err(); err != nil {
		return nil, store.ErrRecordNotFound
	}

	if len(records) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return records, nil
}
