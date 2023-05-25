package teststore

import (
	"employee-service/internal/app/model"
	"employee-service/internal/app/store"
)

type Store struct {
	employeeRepository *EmployeeRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Employee() store.EmployeeRepository {
	if s.employeeRepository != nil {
		return s.employeeRepository
	}

	s.employeeRepository = &EmployeeRepository{
		store:     s,
		employees: make(map[string]*model.Employee),
	}

	return s.employeeRepository
}
