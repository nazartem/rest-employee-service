package store

import "employee-service/internal/app/model"

type EmployeeRepository interface {
	Create(*model.Employee) error
	FindById(int) (*model.Employee, error)
	FindByCompany(int) ([]*model.Employee, error)
	FindByDepartment(int, string) ([]*model.Employee, error)
	PartiallyUpdate(int, *model.Employee) error
	Delete(int) error
}
