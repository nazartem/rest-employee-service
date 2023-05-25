package sqlstore

import (
	"database/sql"
	"employee-service/internal/app/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db                 *sql.DB
	employeeRepository *EmployeeRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Employee() store.EmployeeRepository {
	if s.employeeRepository != nil {
		return s.employeeRepository
	}

	s.employeeRepository = &EmployeeRepository{
		store: s,
	}

	return s.employeeRepository
}
