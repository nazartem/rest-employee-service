package teststore

import (
	"employee-service/internal/app/model"
)

type EmployeeRepository struct {
	store     *Store
	employees map[string]*model.Employee
}

func (r *EmployeeRepository) Create(e *model.Employee) error {
	r.employees[e.Surname] = e
	e.ID = len(r.employees)

	return nil
}

func (r *EmployeeRepository) FindByCompany(companyID int) ([]*model.Employee, error) {
	employess := make([]*model.Employee, 0)

	for surname := range r.employees {
		if r.employees[surname].CompanyID == companyID {
			employess = append(employess, r.employees[surname])
		}
	}

	return employess, nil
}

func (r *EmployeeRepository) Delete(id int) error {
	for surname := range r.employees {
		if r.employees[surname].ID == id {
			delete(r.employees, surname)
		}
	}

	return nil
}
