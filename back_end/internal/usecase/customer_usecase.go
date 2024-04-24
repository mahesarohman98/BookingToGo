package usecase

import (
	"booking_to_go/internal/domain/customer"
	"context"
	"time"

	commonerrors "booking_to_go/internal/errors"

	"github.com/google/uuid"
)

type CustomerUsecase struct {
	repo customer.Repository
}

func NewCustomerUsecase(repo customer.Repository) *CustomerUsecase {
	return &CustomerUsecase{
		repo: repo,
	}
}

func (uc *CustomerUsecase) GetCountries(ctx context.Context) ([]customer.Country, error) {
	return uc.repo.GetCountries(ctx)
}

func (uc *CustomerUsecase) GetListCustomers(ctx context.Context, perPage int, page int) ([]customer.Customer, int, error) {
	return uc.repo.GetListCustomers(ctx, perPage, page)
}

func (uc *CustomerUsecase) GetCustomerByID(ctx context.Context, customerID string) (*customer.Customer, error) {
	return uc.repo.GetCustomerByID(ctx, customerID)
}

type CreateCustomerRequest struct {
	Name        string
	Dob         time.Time
	PhoneNumber string
	Email       string
	CountryID   int
	FamilyList  []FamilyListRequest
}

type FamilyListRequest struct {
	Relation    string
	Name        string
	DateOfBirth time.Time
}

func (uc *CustomerUsecase) CreateCustomer(ctx context.Context, request CreateCustomerRequest) error {

	country, err := uc.repo.GetCountryByID(ctx, request.CountryID)
	if err != nil {
		return err
	}

	if country == nil {
		return commonerrors.NewIncorrectInputError("country not found", "not-found-country")
	}

	newCustomer, err := customer.NewCustomer(uuid.New().String(), request.Name, request.Dob, request.PhoneNumber, request.Email, *country)
	if err != nil {
		return err
	}

	for _, family := range request.FamilyList {
		if err := newCustomer.AddFamilyMember(uuid.New().String(), family.Relation, family.Name, family.DateOfBirth); err != nil {
			return err
		}
	}

	if err := uc.repo.CreateCustomer(ctx, newCustomer); err != nil {
		return err
	}

	return nil
}

type UpdateCustomerRequest struct {
	ID          string
	Name        string
	Dob         *time.Time
	PhoneNumber string
	Email       string
	CountryID   int
	FamilyList  []UpdateFamilyListRequest
}

type UpdateFamilyListRequest struct {
	ID          string
	Relation    string
	Name        string
	DateOfBirth time.Time
}

func (uc *CustomerUsecase) UpdateCustomer(ctx context.Context, request UpdateCustomerRequest) error {

	country, err := uc.repo.GetCountryByID(ctx, request.CountryID)
	if err != nil {
		return err
	}

	if err := uc.repo.UpdateCustomer(
		ctx,
		request.ID,
		func(customer *customer.Customer) (*customer.Customer, error) {

			err := customer.UpdateCustumer(request.Name, request.Dob, request.PhoneNumber, request.Email, country)
			if err != nil {
				return nil, err
			}

			for _, family := range request.FamilyList {
				if family.ID == "" {
					if err := customer.AddFamilyMember(uuid.New().String(), family.Relation, family.Name, family.DateOfBirth); err != nil {
						return nil, err
					}
				} else {
					if err := customer.UpdateFamilyMember(family.ID, family.Relation, family.Name, &family.DateOfBirth); err != nil {
						return nil, err
					}
				}
			}

			return customer, nil
		},
	); err != nil {
		return err
	}

	return nil
}

type AddFamilyMember struct {
	CustomerID string
	List       []FamilyMember
}

type FamilyMember struct {
	Relation    string
	Name        string
	DateOfBirth time.Time
}

func (uc *CustomerUsecase) AddFamilyMember(ctx context.Context, request AddFamilyMember) error {
	if err := uc.repo.UpdateCustomer(
		ctx,
		request.CustomerID,
		func(customer *customer.Customer) (*customer.Customer, error) {

			for _, f := range request.List {
				customer.AddFamilyMember(uuid.New().String(), f.Relation, f.Name, f.DateOfBirth)
			}

			return customer, nil
		},
	); err != nil {
		return err
	}

	return nil
}

type RemoveFamilyMember struct {
	CustomerID string
	FamilyID   string
}

func (uc *CustomerUsecase) RemoveFamilyMember(ctx context.Context, request RemoveFamilyMember) error {

	if err := uc.repo.UpdateCustomer(
		ctx,
		request.CustomerID,
		func(customer *customer.Customer) (*customer.Customer, error) {

			if err := customer.RemoveFamilyMember(request.FamilyID); err != nil {
				return nil, err
			}

			return customer, nil
		},
	); err != nil {
		return err
	}

	return nil
}

type DeleteCustomerRequest struct {
	ID string
}

func (uc *CustomerUsecase) DeleteCustomerByID(ctx context.Context, request DeleteCustomerRequest) error {

	if err := uc.repo.DeleteCustomer(ctx, request.ID); err != nil {
		return err
	}

	return nil
}
