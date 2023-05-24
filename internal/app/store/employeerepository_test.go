package store_test

import (
	"employee-service/internal/app/model"
	"employee-service/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmployeeRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("companies", "departments", "employees", "passports")

	e, err := s.Employee().Create(&model.Employee{
		Name:      "Олег",
		Surname:   "Олегов",
		Phone:     "+7228228",
		CompanyID: 1,
		Passport: struct {
			Type   string
			Number string
		}{Type: "Мультипаспорт", Number: "1488"},
		Department: struct {
			Name  string
			Phone string
		}{Name: "Отдел острых крыльев", Phone: "228-228"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, e)

}
