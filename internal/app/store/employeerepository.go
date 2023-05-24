package store

import (
	"database/sql"
	"employee-service/internal/app/model"
	"log"
	"strings"
)

type EmployeeRepository struct {
	store *Store
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *EmployeeRepository) Create(e *model.Employee) (*model.Employee, error) {
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
			return nil, err
		}
	case err != nil:
		return nil, err
	default:
		log.Printf("Find passportID in passports : %d\n", passportId)
	}

	//// Поиск companyId в таблице companies
	//qCompanyId := `
	//	SELECT
	//		*
	//	FROM
	//	    companies
	//	WHERE id = $1
	//	`
	//if err = r.store.db.QueryRow(
	//	qCompanyId,
	//	e.CompanyID,
	//).Err(); err != nil { // ---- Проверить тут
	//	return nil, err
	//}

	// Поиск departmentId в таблице departments
	qDepartment := `
		SELECT
		    id
		FROM
		    departments
		WHERE
		    name = $1 AND phone = $2
		`
	err = r.store.db.QueryRow(
		qDepartment,
		e.Department.Name,
		e.Department.Phone,
	).Scan(&departmentId)

	switch {
	// Если Id не найден, создаем новую запись
	case err == sql.ErrNoRows:
		qPassportId = `
				INSERT INTO departments 
					(name, phone) 
				VALUES 
					   ($1, $2) 
				RETURNING id
				`
		err = r.store.db.QueryRow(
			qPassportId,
			e.Department.Name,
			e.Department.Phone,
		).Scan(&departmentId)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
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
	if err = r.store.db.QueryRow(
		qEmployee,
		e.Name,
		e.Surname,
		e.Phone,
		e.CompanyID,
		passportId,
		departmentId,
	).Scan(&e.ID); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EmployeeRepository) FindById(id int) (*model.Employee, error) {
	return nil, nil
}
