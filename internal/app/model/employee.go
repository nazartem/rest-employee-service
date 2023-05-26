package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Employee struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Phone     string `json:"phone"`
	CompanyID int    `json:"company_id"`
	Passport  struct {
		Type   string
		Number string
	} `json:"passport"`
	Department struct {
		Name  string
		Phone string
	} `json:"department"`
}

func (e *Employee) Validate() error {
	if err := validation.ValidateStruct(
		e,
		validation.Field(&e.Name, validation.Required, is.UTFLetter),
		validation.Field(&e.Surname, validation.Required, is.UTFLetter),
		validation.Field(&e.Phone, validation.Required, is.Digit),
		validation.Field(&e.CompanyID, validation.Required),

		//validation.Length(6, 100)
	); err != nil {
		return err
	}

	passport := &e.Passport
	if err := validation.ValidateStruct(
		passport,
		validation.Field(&passport.Type, validation.Required, is.UTFLetter),
		validation.Field(&passport.Number, validation.Required, is.Digit),
	); err != nil {
		return err
	}

	department := &e.Department
	if err := validation.ValidateStruct(
		department,
		validation.Field(&department.Name, validation.Required),
		validation.Field(&department.Phone, validation.Required, is.Digit),
	); err != nil {
		return err
	}

	return nil
}
