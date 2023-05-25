package model_test

import (
	"employee-service/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmployee_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		e       func() *model.Employee
		isValid bool
	}{
		{
			name: "valid",
			e: func() *model.Employee {
				return model.TestEmployee(t)
			},
			isValid: true,
		},
		{
			name: "empty name",
			e: func() *model.Employee {
				e := model.TestEmployee(t)
				e.Name = ""

				return e
			},
			isValid: false,
		},
		{
			name: "empty surname",
			e: func() *model.Employee {
				e := model.TestEmployee(t)
				e.Surname = ""

				return e
			},
			isValid: false,
		},
		{
			name: "invalid name",
			e: func() *model.Employee {
				e := model.TestEmployee(t)
				e.Name = "1234"

				return e
			},
			isValid: false,
		},
		{
			name: "empty phone",
			e: func() *model.Employee {
				e := model.TestEmployee(t)
				e.Phone = ""

				return e
			},
			isValid: false,
		},
		{
			name: "empty passport",
			e: func() *model.Employee {
				e := model.TestEmployee(t)
				e.Passport = struct {
					Type   string
					Number string
				}{}

				return e
			},
			isValid: false,
		},
		{
			name: "empty department",
			e: func() *model.Employee {
				e := model.TestEmployee(t)
				e.Department = struct {
					Name  string
					Phone string
				}{}

				return e
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.e().Validate())
			} else {
				assert.Error(t, tc.e().Validate())
			}
		})
	}
}
