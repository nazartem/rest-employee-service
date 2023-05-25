package model

type Employee struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Phone     string `json:"phone"`
	CompanyID int    `json:"companyID"`
	Passport  struct {
		Type   string
		Number string
	} `json:"passport"`
	Department struct {
		Name  string
		Phone string
	} `json:"department"`
}
