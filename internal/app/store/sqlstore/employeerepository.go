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

	var passportId, departmentId int

	// Поиск Id паспорта в таблице passports
	qPassportId := `
		SELECT
		    id
		FROM
		    passports
		WHERE
		    type = $1 AND number = $2
		`
	err := r.store.db.QueryRow(
		qPassportId,
		e.Passport.Type,
		e.Passport.Number,
	).Scan(&passportId)

	switch {
	// Если Id не найден, создаем новую запись
	case err == sql.ErrNoRows:
		qPassportId = `
				INSERT INTO passports 
					(type, number) 
				VALUES 
					   ($1, $2) 
				RETURNING id
				`
		err = r.store.db.QueryRow(
			qPassportId,
			e.Passport.Type,
			e.Passport.Number,
		).Scan(&passportId)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		log.Printf("Find passportID in passports : %d\n", passportId)
	}

	// Поиск departmentId в таблице departments
	qDepartmentId := `
		SELECT
		    id
		FROM
		    departments
		WHERE
		    name = $1 AND phone = $2
		`
	err = r.store.db.QueryRow(
		qDepartmentId,
		e.Department.Name,
		e.Department.Phone,
	).Scan(&departmentId)

	switch {
	// Если Id не найден, создаем новую запись
	case err == sql.ErrNoRows:
		qDepartmentId = `
				INSERT INTO departments 
					(name, phone) 
				VALUES 
					   ($1, $2) 
				RETURNING id
				`
		err = r.store.db.QueryRow(
			qDepartmentId,
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
					(name, surname, phone, companyid, passportid, departmentid) 
				VALUES 
					   ($1, $2, $3, $4, $5, $6) 
				RETURNING id
				`
	return r.store.db.QueryRow(
		qEmployee,
		e.Name,
		e.Surname,
		e.Phone,
		e.CompanyID,
		passportId,
		departmentId,
	).Scan(&e.ID)
}

func (r *EmployeeRepository) FindByCompany(companyID int) ([]*model.Employee, error) {
	qEmployee := `
		SELECT
		    id, name, surname, phone, passportid, departmentid
		FROM
		    employees
		WHERE
		    companyid = $1
		`

	rows, err := r.store.db.Query(qEmployee, companyID)
	if err != nil {
		return nil, err
	}

	employees := make([]*model.Employee, 0)

	for rows.Next() {
		e := &model.Employee{}
		var passportId, departmentId int

		err = rows.Scan(
			&e.ID,
			&e.Name,
			&e.Surname,
			&e.Phone,
			&passportId,
			&departmentId,
		)
		if err != nil {
			return nil, err
		}

		// Поиск type, number по passportId в таблице passports
		qPassport := `
			SELECT
				type, number
			FROM
				passports
			WHERE
				id = $1
			`
		if err = r.store.db.QueryRow(
			qPassport,
			passportId,
		).Scan(
			&e.Passport.Type,
			&e.Passport.Number,
		); err != nil {
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
		//r.logger.Err.Println(newErr)
		return newErr
	}

	return nil
}
