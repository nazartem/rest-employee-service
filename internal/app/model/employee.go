package model

type Employee struct {
	ID        int
	Name      string
	Surname   string
	Phone     string
	CompanyID int
	Passport  struct {
		Type   string
		Number string
	}
	Department struct {
		Name  string
		Phone string
	}
}
