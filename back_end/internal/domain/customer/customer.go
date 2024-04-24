package customer

import (
	commonerrors "booking_to_go/internal/errors"
	"net/mail"
	"time"

	"github.com/nyaruka/phonenumbers"
)

type Customer struct {
	ID string

	Name        string
	DateOfBirth time.Time
	PhoneNumber string
	Email       string

	Nationality Country

	FamilyList []FamilyList
}

var (
	ErrInvalidEmail       = commonerrors.NewIncorrectInputError("email not valid", "validation-error")
	ErrInvalidPhoneNumber = commonerrors.NewIncorrectInputError("phone number not valid", "validation-error")
)

func validate(email string, phoneNumber string, defaultCountry Country) error {
	if email != "" {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return ErrInvalidEmail
		}
	}

	if phoneNumber != "" {
		_, err := phonenumbers.Parse(phoneNumber, defaultCountry.Code)
		if err != nil {
			return ErrInvalidPhoneNumber
		}
	}

	return nil
}

func (c *Customer) UpdateCustumer(name string, dob *time.Time, phoneNumber string, email string, nationality *Country) error {

	if name != "" {
		c.Name = name
	}

	if dob != nil {
		c.DateOfBirth = *dob
	}

	if phoneNumber != "" {
		c.PhoneNumber = phoneNumber
	}

	if email != "" {
		c.Email = email
	}

	if nationality != nil {
		c.Nationality = *nationality
	}

	if err := validate(c.Email, c.PhoneNumber, c.Nationality); err != nil {
		return err
	}

	return nil
}

func NewCustomer(id string, name string, dob time.Time, phoneNumber string, email string, nationality Country) (*Customer, error) {

	if err := validate(email, phoneNumber, nationality); err != nil {
		return &Customer{}, err
	}

	return &Customer{
		ID:          id,
		Name:        name,
		DateOfBirth: dob,
		PhoneNumber: phoneNumber,
		Email:       email,
		Nationality: nationality,
		FamilyList:  []FamilyList{},
	}, nil
}
