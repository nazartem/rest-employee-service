package store

import "employee-service/internal/app/model"

type EmployeeRepository interface {
	Create(*model.Employee) error
	FindByCompany(int) ([]*model.Employee, error)
	Delete(int) error
}
