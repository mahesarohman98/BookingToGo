package customer

import "context"

type Repository interface {
	GetCountries(ctx context.Context) ([]Country, error)
	GetCountryByID(ctx context.Context, countryID int) (*Country, error)

	GetCustomerByID(ctx context.Context, id string) (*Customer, error)
	GetListCustomers(ctx context.Context, perPage int, page int) ([]Customer, int, error)

	CreateCustomer(
		ctx context.Context,
		customer *Customer,
	) error

	UpdateCustomer(
		ctx context.Context,
		customerID string,
		Fn func(customer *Customer) (*Customer, error),
	) error

	DeleteCustomer(
		ctx context.Context,
		customerID string,
	) error
}
