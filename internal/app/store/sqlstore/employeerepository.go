package sqlstore

import (
	"database/sql"
	"employee-service/internal/app/model"
	"errors"
	"log"
)

type EmployeeRepository struct {
	store *Store
}

func (r *EmployeeRepository) Create(e *model.Employee) error {
	if err := e.Validate(); err != nil {
		return err
	}

	var departmentId int

	// Поиск departmentId в таблице departments
	qDepartmentId := `
		SELECT
		    id
		FROM
		    departments
		WHERE
		    company_id = $1 AND name = $2 AND phone = $3
		`
	err := r.store.db.QueryRow(
		qDepartmentId,
		e.CompanyID,
		e.Department.Name,
		e.Department.Phone,
	).Scan(&departmentId)

	switch {
	// Если Id не найден, создаем новую запись
	case err == sql.ErrNoRows:
		qDepartmentId = `
				INSERT INTO departments 
					(company_id, name, phone) 
				VALUES 
					   ($1, $2, $3) 
				RETURNING id
				`
		err = r.store.db.QueryRow(
			qDepartmentId,
			e.CompanyID,
			e.Department.Name,
			e.Department.Phone,
		).Scan(&departmentId)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		log.Printf("Find departmentId in departments : %d\n", departmentId)
	}

	// Создание новой записи в таблице employees
	qEmployee := `
				INSERT INTO employees 
					(name, surname, phone, company_id, passport_type, passport_number, department_id) 
				VALUES 
					   ($1, $2, $3, $4, $5, $6, $7) 
				RETURNING id
				`
	return r.store.db.QueryRow(
		qEmployee,
		e.Name,
		e.Surname,
		e.Phone,
		e.CompanyID,
		e.Passport.Type,
		e.Passport.Number,
		departmentId,
	).Scan(&e.ID)
}

func (r *EmployeeRepository) FindByCompany(companyID int) ([]*model.Employee, error) {
	qEmployee := `
		SELECT
		    id, name, surname, phone, passport_type, passport_number, department_id
		FROM
		    employees
		WHERE
		    company_id = $1
		`

	rows, err := r.store.db.Query(qEmployee, companyID)
	if err != nil {
		return nil, err
	}

	employees := make([]*model.Employee, 0)

	for rows.Next() {
		e := &model.Employee{}
		var departmentId int

		err = rows.Scan(
			&e.ID,
			&e.Name,
			&e.Surname,
			&e.Phone,
			&e.Passport.Type,
			&e.Passport.Number,
			&departmentId,
		)
		if err != nil {
			return nil, err
		}

		// Поиск name, phone по departmentId в таблице departments
		qDepartment := `
			SELECT
				name, phone
			FROM
				departments
			WHERE
				id = $1
			`
		if err = r.store.db.QueryRow(
			qDepartment,
			departmentId,
		).Scan(
			&e.Department.Name,
			&e.Department.Phone,
		); err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *EmployeeRepository) FindByDepartment(companyID int, department string) ([]*model.Employee, error) {
	var departmentId int
	var departmentPhone string

	qDepartmentId := `
		SELECT
			id, phone
		FROM
			departments
		WHERE
		company_id = $1 AND name = $2
		`

	if err := r.store.db.QueryRow(
		qDepartmentId,
		companyID,
		department,
	).Scan(&departmentId, &departmentPhone); err != nil {
		return nil, err
	}

	qEmployee := `
		SELECT
		    id, name, surname, phone, passport_type, passport_number
		FROM
		    employees
		WHERE
		    company_id = $1 AND department_id = $2
		`

	rows, err := r.store.db.Query(qEmployee, companyID, departmentId)
	if err != nil {
		return nil, err
	}

	employees := make([]*model.Employee, 0)

	for rows.Next() {
		e := &model.Employee{
			CompanyID: companyID,
			Department: struct {
				Name  string
				Phone string
			}{Name: department, Phone: departmentPhone},
		}

		err = rows.Scan(
			&e.ID,
			&e.Name,
			&e.Surname,
			&e.Phone,
			&e.Passport.Type,
			&e.Passport.Number,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *EmployeeRepository) FindById(id int) (*model.Employee, error) {
	var departmentId int
	e := &model.Employee{}

	qEmployee := `
		SELECT
		    name, surname, phone, company_id, passport_type, passport_number, department_id
		FROM
		    employees
		WHERE
		    id = $1
		`
	if err := r.store.db.QueryRow(
		qEmployee,
		id,
	).Scan(
		&e.Name,
		&e.Surname,
		&e.Phone,
		&e.CompanyID,
		&e.Passport.Type,
		&e.Passport.Number,
		&departmentId,
	); err != nil {
		return nil, err
	}

	qDepartmentId := `
		SELECT
			name, phone
		FROM
			departments
		WHERE
		id = $1
		`

	if err := r.store.db.QueryRow(
		qDepartmentId,
		departmentId,
	).Scan(&e.Department.Name, &e.Department.Phone); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EmployeeRepository) PartiallyUpdate(id int, e *model.Employee) error {
	var departmentId int

	// Поиск departmentId в таблице departments
	qDepartmentId := `
		SELECT
		    id
		FROM
		    departments
		WHERE
		    company_id = $1 AND name = $2 AND phone = $3
		`
	err := r.store.db.QueryRow(
		qDepartmentId,
		e.CompanyID,
		e.Department.Name,
		e.Department.Phone,
	).Scan(&departmentId)

	switch {
	// Если Id не найден, создаем новую запись
	case err == sql.ErrNoRows:
		qDepartmentId = `
				INSERT INTO departments 
					(company_id, name, phone) 
				VALUES 
					   ($1, $2, $3) 
				RETURNING id
				`
		err = r.store.db.QueryRow(
			qDepartmentId,
			e.CompanyID,
			e.Department.Name,
			e.Department.Phone,
		).Scan(&departmentId)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		log.Printf("Find departmentId in departments : %d\n", departmentId)
	}

	// Обновление принятых полей в employees
	qUpdateEmployee := `
		UPDATE
			employees
		SET
		    name = $1,
		    surname = $2,
		    phone = $3,
		    company_id = $4,
		    passport_type = $5,
		    passport_number = $6,
			department_id = $7
		WHERE
		    id = $8
		`

	//args := make([]interface{}, 0)

	//if employee.Name != "" {
	//	args = append(args, employee.Name)
	//	str := fmt.Sprintf(" name=$%d,", len(args))
	//	updateQuery += str
	//}
	//if employee.Surname != "" {
	//	args = append(args, employee.Surname)
	//	str := fmt.Sprintf(" surname=$%d,", len(args))
	//	updateQuery += str
	//}
	//if employee.Phone != "" {
	//	args = append(args, employee.Phone)
	//	str := fmt.Sprintf(" phone=$%d,", len(args))
	//	updateQuery += str
	//}
	//if employee.Passport.Type != "" {
	//	args = append(args, employee.Passport.Type)
	//	str := fmt.Sprintf(" passport_type=$%d,", len(args))
	//	updateQuery += str
	//}
	//if employee.Passport.Number != "" {
	//	args = append(args, employee.Passport.Number)
	//	str := fmt.Sprintf(" passport_number=$%d,", len(args))
	//	updateQuery += str
	//}
	//if len(args) == 0 {
	//	// если запрос UPDATE не содержит данных для обновления
	//	newErr := errors.New("no data to update")
	//	return newErr
	//}
	//args = append(args, id)
	//str := fmt.Sprintf(" WHERE id=$%d", len(args))
	//updateQuery = strings.TrimSuffix(updateQuery, ",") + str
	//
	//_, err := r.store.db.Exec(updateQuery, args...)
	//if err != nil {
	//	return err
	//}

	err = r.store.db.QueryRow(
		qUpdateEmployee,
		e.Name,
		e.Surname,
		e.Phone,
		e.CompanyID,
		e.Passport.Type,
		e.Passport.Number,
		departmentId,
		id,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) Delete(id int) error {
	q := `DELETE FROM employees WHERE id = $1`
	commandTag, err := r.store.db.Exec(q, id)
	if err != nil {
		return err
	}

	rows, _ := commandTag.RowsAffected()
	if rows != 1 {
		newErr := errors.New("no row found to delete")
		// TODO logging
		log.Printf(newErr.Error())
		return newErr
	}

	return nil
}
