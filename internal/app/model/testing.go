package model

import "testing"

var (
	names            = []string{"Андрей", "Борис", "Виталий", "Григорий", "Дмитрий", "Евгений"}
	surnames         = []string{"Андрев", "Борисов", "Витальев", "Григорьев", "Дмитриев", "Евгеньев"}
	phones           = []string{"1232", "4563", "7891", "1344", "2568", "9760"}
	types            = []string{"4563", "1232", "1344", "7891", "9760", "2568"}
	numbers          = []string{"7891", "1344", "9760", "2568", "1232", "4563"}
	departmentNames  = []string{"Отдел №1", "Отдел №2", "Отдел №3", "Отдел №4", "Отдел №5", "Отдел №6"}
	departmentPhones = []string{"11-11", "22-22", "33-33", "44-44", "55-55", "66-66"}
)

func TestEmployee(t *testing.T) *Employee {
	return &Employee{
		Name:      "Корбен",
		Surname:   "Даллас",
		Phone:     "+1234567",
		CompanyID: 1,
		Passport: struct {
			Type   string
			Number string
		}{Type: "Мультипаспорт", Number: "2345"},
		Department: struct {
			Name  string
			Phone string
		}{Name: "Отдел №33", Phone: "65-65-65"},
	}
}

func TestArrEmployees(t *testing.T) []*Employee {
	employeesArr := make([]*Employee, 0)

	for i := 0; i < len(names); i++ {
		e := &Employee{
			Name:      names[i],
			Surname:   surnames[i],
			Phone:     phones[i],
			CompanyID: 1,
			Passport: struct {
				Type   string
				Number string
			}{Type: types[i], Number: numbers[i]},
			Department: struct {
				Name  string
				Phone string
			}{Name: departmentNames[i], Phone: departmentPhones[i]},
		}

		employeesArr = append(employeesArr, e)
	}

	return employeesArr
}
