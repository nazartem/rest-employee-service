package sqlstore_test

import (
	"employee-service/internal/app/model"
	"employee-service/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmployeeRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("companies", "departments", "employees", "passports")

	s := sqlstore.New(db)
	e := model.TestEmployee(t)
	assert.NoError(t, s.Employee().Create(e))
	assert.NotNil(t, e)
}

func TestEmployeeRepository_FindByCompany(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("companies", "departments", "employees", "passports")

	s := sqlstore.New(db)
	companyID := 1
	e, err := s.Employee().FindByCompany(companyID)
	assert.Nil(t, err)

	testEmployees := model.TestArrEmployees(t)
	for i := 0; i < len(testEmployees); i++ {
		s.Employee().Create(testEmployees[i])
	}

	e, err = s.Employee().FindByCompany(companyID)
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, len(testEmployees), len(e))
}
