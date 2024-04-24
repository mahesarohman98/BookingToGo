package repositories

import (
	"booking_to_go/internal/domain/customer"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type CustomerPGRepository struct {
	db *sqlx.DB
}

func NewCustomerPGRepository(db *sqlx.DB) *CustomerPGRepository {
	return &CustomerPGRepository{
		db: db,
	}
}

// const postgresDeadlockErrorCode = "40P01"
// const postgresViolatedConstrainKey = "23503"

type pgNationality struct {
	ID   int    `db:"nationality_id"`
	Name string `db:"nationality_name"`
	Code string `db:"nationality_code"`
}

type pgCustomer struct {
	ID          string    `db:"cst_id"`
	NationalID  int       `db:"nationality_id"`
	Name        string    `db:"cst_name"`
	DateOfBirth time.Time `db:"cst_dob"`
	PhoneNumber string    `db:"cst_phone_num"`
	Email       string    `db:"cst_email"`
}

type pgCustomerAggregate struct {
	CustomerID  string    `db:"cst_id"`
	Name        string    `db:"cst_name"`
	DateOfBirth time.Time `db:"cst_dob"`
	PhoneNumber string    `db:"cst_phone_num"`
	Email       string    `db:"cst_email"`

	NationalID   int    `db:"nationality_id"`
	NationalName string `db:"nationality_name"`
	NationalCode string `db:"nationality_code"`

	FamilyMemberID          *string    `db:"fl_id"`
	FamilyMemberRelation    *string    `db:"fl_relation"`
	FamilyMemberName        *string    `db:"fl_name"`
	FamilyMemberDateOfBirth *time.Time `db:"fl_dob"`
}

func unMarshallCustomerFromDB(dbCustomers []pgCustomerAggregate) (*customer.Customer, error) {

	var customerAggr *customer.Customer
	for i, dbCustomer := range dbCustomers {
		if i == 0 {
			country := customer.NewContry(dbCustomer.NationalID, dbCustomer.NationalName, dbCustomer.NationalCode)

			c, err := customer.NewCustomer(dbCustomer.CustomerID, dbCustomer.Name, dbCustomer.DateOfBirth, dbCustomer.PhoneNumber, dbCustomer.Email, *country)
			if err != nil {
				return nil, err
			}

			customerAggr = c
		}

		if dbCustomer.FamilyMemberID != nil {
			customerAggr.AddFamilyMember(*dbCustomer.FamilyMemberID, *dbCustomer.FamilyMemberRelation, *dbCustomer.FamilyMemberName, *dbCustomer.FamilyMemberDateOfBirth)
		}

	}

	return customerAggr, nil
}

func (repo *CustomerPGRepository) GetCountries(ctx context.Context) ([]customer.Country, error) {
	countries := make([]customer.Country, 0)

	query := `SELECT * FROM nationality`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return []customer.Country{}, err
	}

	for rows.Next() {
		n := pgNationality{}
		err := rows.Scan(&n.ID, &n.Name, &n.Code)
		if err != nil {
			return []customer.Country{}, err
		}

		country := customer.NewContry(n.ID, n.Name, n.Code)

		countries = append(countries, *country)
	}

	return countries, nil
}
func (repo *CustomerPGRepository) GetCountryByID(ctx context.Context, countryID int) (*customer.Country, error) {
	nationality := pgNationality{}

	query := "SELECT * FROM nationality n WHERE n.nationality_id = $1"

	err := repo.db.GetContext(ctx, &nationality, query, countryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get nationality from db")
	}

	return customer.NewContry(nationality.ID, nationality.Name, nationality.Code), nil

}

func (repo *CustomerPGRepository) CreateCustomer(
	ctx context.Context,
	customer *customer.Customer,
) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	defer func() {
		err = repo.finishTransaction(err, tx)
	}()

	err = repo.insertCustomer(tx, customer)
	if err != nil {
		return errors.Wrap(err, "unable to insert customer")
	}

	err = repo.batchInsertFamilyMember(tx, customer)
	if err != nil {
		return errors.Wrap(err, "unable to bulk insert family member")
	}

	return nil
}

func (repo *CustomerPGRepository) updateCustomer(db sqlContextGetter, customer *customer.Customer) error {
	pgCustomer := pgCustomer{
		ID:          customer.ID,
		NationalID:  customer.Nationality.ID,
		Name:        customer.Name,
		DateOfBirth: customer.DateOfBirth,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}

	_, err := db.NamedExec(`
		UPDATE Customers SET 
			nationality_id = :nationality_id, 
			cst_name = :cst_name,
			cst_dob = :cst_dob,
			cst_phone_num = :cst_phone_num,
			cst_email = :cst_email
		WHERE cst_id = :cst_id
	`, pgCustomer)

	return err
}

func (repo *CustomerPGRepository) insertCustomer(db sqlContextGetter, customer *customer.Customer) error {
	pgCustomer := pgCustomer{
		ID:          customer.ID,
		NationalID:  customer.Nationality.ID,
		Name:        customer.Name,
		DateOfBirth: customer.DateOfBirth,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}

	_, err := db.NamedExec(
		`INSERT INTO 
			Customers (cst_id, nationality_id, cst_name, cst_dob, cst_phone_num, cst_email) 
		VALUES 
			(:cst_id, :nationality_id, :cst_name, :cst_dob, :cst_phone_num, :cst_email)`, pgCustomer)

	return err
}

func (repo *CustomerPGRepository) batchInsertFamilyMember(db sqlContextGetter, customer *customer.Customer) error {
	if len(customer.FamilyList) <= 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(customer.FamilyList))
	valueArgs := make([]interface{}, 0, len(customer.FamilyList)*5)

	i := 1
	for _, family := range customer.FamilyList {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i, i+1, i+2, i+3, i+4))
		valueArgs = append(valueArgs, family.ID)
		valueArgs = append(valueArgs, customer.ID)
		valueArgs = append(valueArgs, family.Relation)
		valueArgs = append(valueArgs, family.Name)
		valueArgs = append(valueArgs, family.DateOfBirth)
		i += 5
	}

	query := fmt.Sprintf(`
		INSERT INTO FamilyLists AS c
			(fl_id, cst_id, fl_relation, fl_name, fl_dob)
		VALUES %s
		ON CONFLICT (fl_id) DO UPDATE SET
			fl_relation = EXCLUDED.fl_relation,  
			fl_name = EXCLUDED.fl_name,  
			fl_dob = EXCLUDED.fl_dob  
		WHERE c.cst_id = $%d
			`,
		strings.Join(valueStrings, ","), i,
	)

	valueArgs = append(valueArgs, customer.ID)

	_, err := db.Exec(query, valueArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to bulk insert family member")
	}

	return nil

}

func (repo *CustomerPGRepository) GetListCustomers(ctx context.Context, perPage int, page int) ([]customer.Customer, int, error) {
	customers := make([]customer.Customer, 0)
	var total int
	query := `
		SELECT *, count(*) OVER() AS full_count  FROM Customers 
			JOIN nationality USING (nationality_id) 
		LIMIT $1 OFFSET $2
		`

	rows, err := repo.db.QueryContext(ctx, query, perPage, page)
	if err != nil {
		return []customer.Customer{}, 0, err
	}

	for rows.Next() {
		c := pgCustomerAggregate{}

		err := rows.Scan(&c.NationalID, &c.CustomerID, &c.Name, &c.DateOfBirth, &c.PhoneNumber, &c.Email, &c.NationalName, &c.NationalCode, &total)
		if err != nil {
			return []customer.Customer{}, 0, err
		}

		country := customer.NewContry(c.NationalID, c.NationalName, c.NationalCode)

		customers = append(customers, customer.Customer{
			ID:          c.CustomerID,
			Name:        c.Name,
			DateOfBirth: c.DateOfBirth,
			PhoneNumber: c.PhoneNumber,
			Email:       c.Email,
			Nationality: *country,
			FamilyList:  []customer.FamilyList{},
		})

	}

	return customers, total, nil
}

func (repo *CustomerPGRepository) GetCustomerByID(ctx context.Context, id string) (*customer.Customer, error) {
	customer, err := repo.getCustomerByID(ctx, repo.db, id)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, sql.ErrNoRows
	}

	return customer, err
}

// sqlContextGetter is an interface provided both by transaction and standard db connection
type sqlContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

func (repo *CustomerPGRepository) getCustomerByID(
	ctx context.Context,
	db sqlContextGetter,
	id string,
) (*customer.Customer, error) {
	dbCustomer := []pgCustomerAggregate{}

	query := `
		SELECT * FROM Customers
			JOIN nationality using (nationality_id)
			LEFT JOIN FamilyLists using(cst_id)
		WHERE cst_id = $1
		`

	rows, err := db.QueryContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := pgCustomerAggregate{}
		err := rows.Scan(&c.CustomerID, &c.NationalID, &c.Name, &c.DateOfBirth, &c.PhoneNumber, &c.Email, &c.NationalName, &c.NationalCode, &c.FamilyMemberID, &c.FamilyMemberRelation, &c.FamilyMemberName, &c.FamilyMemberDateOfBirth)
		if err != nil {
			return nil, err
		}

		dbCustomer = append(dbCustomer, c)
	}

	customer, err := unMarshallCustomerFromDB(dbCustomer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (repo *CustomerPGRepository) UpdateCustomer(
	ctx context.Context,
	customerID string,
	Fn func(customer *customer.Customer) (*customer.Customer, error),
) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	defer func() {
		err = repo.finishTransaction(err, tx)
	}()

	existingCustomer, err := repo.getCustomerByID(ctx, tx, customerID)
	if err != nil {
		return errors.Wrap(err, "unable to get customer")
	}

	customer, err := Fn(existingCustomer)
	if err != nil {
		return err
	}
	err = repo.updateCustomer(tx, customer)
	if err != nil {
		return errors.Wrap(err, "unable to update customer")
	}

	err = repo.deleteFamilyMemberByCustomnerID(tx, customer.ID)
	if err != nil {
		return errors.Wrap(err, "unable to delete family member by customer id")
	}

	err = repo.batchInsertFamilyMember(tx, customer)
	if err != nil {
		return errors.Wrap(err, "unable to bulk insert family member")
	}

	return nil
}

func (repo *CustomerPGRepository) deleteFamilyMemberByCustomnerID(db sqlContextGetter, customerID string) error {
	_, err := db.Exec(`DELETE FROM FamilyLists WHERE cst_id = $1`, customerID)
	return err
}

func (repo *CustomerPGRepository) DeleteCustomer(
	ctx context.Context,
	customerID string,
) error {
	_, err := repo.db.Exec(`DELETE FROM Customers WHERE cst_id = $1`, customerID)
	return err
}

func (repo *CustomerPGRepository) finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}
