package teststore

import (
	"employee-service/internal/app/model"
	"fmt"
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
	employees := make([]*model.Employee, 0)

	for surname := range r.employees {
		if r.employees[surname].CompanyID == companyID {
			employees = append(employees, r.employees[surname])
		}
	}

	return employees, nil
}

func (r *EmployeeRepository) FindByDepartment(companyID int, department string) ([]*model.Employee, error) {
	employees := make([]*model.Employee, 0)

	for surname := range r.employees {
		if r.employees[surname].CompanyID == companyID && r.employees[surname].Department.Name == department {
			employees = append(employees, r.employees[surname])
		}
	}

	return employees, nil
}

func (r *EmployeeRepository) FindById(id int) (*model.Employee, error) {
	for surname := range r.employees {
		if r.employees[surname].ID == id {
			return r.employees[surname], nil
		}
	}

	return nil, fmt.Errorf("no content")
}

func (r *EmployeeRepository) PartiallyUpdate(id int, e *model.Employee) error {
	for surname := range r.employees {
		if r.employees[surname].ID == id {
			r.employees[surname] = e
			return nil
		}
	}

	return fmt.Errorf("no content")
}

func (r *EmployeeRepository) Delete(id int) error {
	for surname := range r.employees {
		if r.employees[surname].ID == id {
			delete(r.employees, surname)
		}
	}

	return nil
}
