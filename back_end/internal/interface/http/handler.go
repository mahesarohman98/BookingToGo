package http

import (
	"booking_to_go/internal/domain/customer"
	"booking_to_go/internal/interface/util/pagination"
	"booking_to_go/internal/usecase"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	usecase *usecase.CustomerUsecase
}

func NewCustomerHandler(usecase *usecase.CustomerUsecase) *CustomerHandler {

	return &CustomerHandler{usecase: usecase}
}

type CreateCustomerRequest struct {
	Name        string               `json:"name"`
	Dob         string               `json:"dob"`
	DateOfBirth time.Time            `json:"-"`
	PhoneNumber string               `json:"phone_number"`
	Email       string               `json:"email"`
	CountryID   int                  `json:"country_id"`
	FamilyList  []*FamilyListRequest `json:"familyList"`
}

type FamilyListRequest struct {
	Relation    string    `json:"relation"`
	Name        string    `json:"name"`
	Dob         string    `json:"dob"`
	DateOfBirth time.Time `json:"-"`
}

func (req *CreateCustomerRequest) Validate() error {

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	v, err := time.Parse("2006-01-02", req.Dob)
	if req.Dob == "" || err != nil {
		return fmt.Errorf("dob is required")
	}
	req.DateOfBirth = v

	if req.PhoneNumber == "" {
		return fmt.Errorf("phone_number is required")
	}

	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	if req.CountryID <= 0 {
		return fmt.Errorf("country_id is required")
	}

	for _, f := range req.FamilyList {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (req *FamilyListRequest) Validate() error {
	if req.Relation == "" {
		return fmt.Errorf("relation is required")
	}

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	v, err := time.Parse("2006-01-02", req.Dob)
	if req.Dob == "" || err != nil {
		return fmt.Errorf("dob is required")
	}
	req.DateOfBirth = v

	return nil
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	req := &CreateCustomerRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		BadRequest("invalid-request", err, w, r)
		return
	}

	// Validate the request data
	if err := req.Validate(); err != nil {
		BadRequest("validation-error", err, w, r)
		return
	}

	createCustomerReq := usecase.CreateCustomerRequest{
		Name:        req.Name,
		Dob:         req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		CountryID:   req.CountryID,
		FamilyList:  []usecase.FamilyListRequest{},
	}

	for _, f := range req.FamilyList {
		createCustomerReq.FamilyList = append(createCustomerReq.FamilyList, usecase.FamilyListRequest{
			Relation:    f.Relation,
			Name:        f.Name,
			DateOfBirth: f.DateOfBirth,
		})
	}

	if err := h.usecase.CreateCustomer(r.Context(), createCustomerReq); err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	writeJSON(w, http.StatusCreated)

}

func (h *CustomerHandler) GetCountries(w http.ResponseWriter, r *http.Request) {

	countries, err := h.usecase.GetCountries(r.Context())
	if err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	writeJSONWithPayload(w, countries)

}

type ListCustomersRequest struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}

type Customer struct {
	ID string `json:"id"`

	Name        string `json:"name"`
	DateOfBirth string `json:"dob"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`

	Nationality string `json:"nationality"`
}

type ListCustomerResponse struct {
	Customers []Customer            `json:"customers"`
	Meta      pagination.Pagination `json:"meta"`
}

func NewListCustomersResponse(customers []customer.Customer, pagination pagination.Pagination) ListCustomerResponse {
	res := ListCustomerResponse{
		Meta:      pagination,
		Customers: []Customer{},
	}

	for _, c := range customers {
		res.Customers = append(res.Customers, Customer{
			ID:          c.ID,
			Name:        c.Name,
			DateOfBirth: c.DateOfBirth.Format("2006-01-02"),
			PhoneNumber: c.PhoneNumber,
			Email:       c.Email,
			Nationality: c.Nationality.String(),
		})
	}

	return res
}

func (h *CustomerHandler) GetListCustomers(w http.ResponseWriter, r *http.Request) {

	req := ListCustomersRequest{}

	var err error

	req.PerPage, err = strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		req.PerPage = 10
		err = nil
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		err = nil
	}

	req.CurrentPage = (page - 1) * req.PerPage
	if page == 0 {
		req.CurrentPage = 1
	}

	customers, total, err := h.usecase.GetListCustomers(r.Context(), req.PerPage, req.CurrentPage)
	if err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	pagination := pagination.NewPagination(req.PerPage, req.CurrentPage, total)

	response := NewListCustomersResponse(customers, pagination)

	writeJSONWithPayload(w, response)

}

type GetCustomerByIDResponse struct {
	ID string `json:"id"`

	Name        string `json:"name"`
	DateOfBirth string `json:"dob"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`

	NationalityID int    `json:"nationality_id"`
	Nationality   string `json:"nationality"`

	FamilyList []FamilyList `json:"family_list"`
}

type FamilyList struct {
	ID          string `json:"id"`
	Relation    string `json:"relation"`
	Name        string `json:"name"`
	DateOfBirth string `json:"dob"`
}

func NewGetCustomerByIDResponse(customer *customer.Customer) GetCustomerByIDResponse {

	if customer == nil {
		return GetCustomerByIDResponse{}
	}

	req := GetCustomerByIDResponse{
		ID:            customer.ID,
		Name:          customer.Name,
		DateOfBirth:   customer.DateOfBirth.Format("2006-01-02"),
		PhoneNumber:   customer.PhoneNumber,
		Email:         customer.Email,
		NationalityID: customer.Nationality.ID,
		Nationality:   customer.Nationality.String(),
		FamilyList:    []FamilyList{},
	}

	for _, c := range customer.FamilyList {
		req.FamilyList = append(req.FamilyList, FamilyList{
			ID:          c.ID,
			Relation:    c.Relation,
			Name:        c.Name,
			DateOfBirth: c.DateOfBirth.Format("2006-01-02"),
		})
	}

	return req
}

func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer, err := h.usecase.GetCustomerByID(r.Context(), vars["id"])
	if err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	writeJSONWithPayload(w, NewGetCustomerByIDResponse(customer))
}

type UpdateCustomerByIDRequest struct {
	ID string `json:"id"`

	Name        string     `json:"name"`
	Dob         string     `json:"dob"`
	DateOfBirth *time.Time `json:"doobs"`
	PhoneNumber string     `json:"phone_number"`
	Email       string     `json:"email"`

	NationalityID int `json:"nationality_id"`

	FamilyList []*FamilyListReq `json:"family_list"`
}

type FamilyListReq struct {
	ID          string    `json:"id"`
	Relation    string    `json:"relation"`
	Name        string    `json:"name"`
	Dob         string    `json:"dob"`
	DateOfBirth time.Time `json:"doobs"`
}

func (req *FamilyListReq) Validate() error {
	if req.Relation == "" {
		return fmt.Errorf("relation is required")
	}

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	v, err := time.Parse("2006-01-02", req.Dob)
	if req.Dob == "" || err != nil {
		return fmt.Errorf("dob is required")
	}
	req.DateOfBirth = v

	return nil

}

func (req *UpdateCustomerByIDRequest) Validate() error {
	if req.ID == "" {
		return fmt.Errorf("dob is required")
	}

	v, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		err = nil
	} else {
		req.DateOfBirth = &v
	}

	for _, f := range req.FamilyList {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil

}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

	req := &UpdateCustomerByIDRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		BadRequest("invalid-request", err, w, r)
		return
	}

	vars := mux.Vars(r)
	req.ID = vars["id"]

	// Validate the request data
	if err := req.Validate(); err != nil {
		BadRequest("validation-error", err, w, r)
		return
	}

	updateCustomerReq := usecase.UpdateCustomerRequest{
		ID:          req.ID,
		Name:        req.Name,
		Dob:         req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		CountryID:   req.NationalityID,
		FamilyList:  []usecase.UpdateFamilyListRequest{},
	}

	for _, f := range req.FamilyList {
		updateCustomerReq.FamilyList = append(updateCustomerReq.FamilyList, usecase.UpdateFamilyListRequest{
			ID:          f.ID,
			Relation:    f.Relation,
			Name:        f.Name,
			DateOfBirth: f.DateOfBirth,
		})
	}

	if err := h.usecase.UpdateCustomer(r.Context(), updateCustomerReq); err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	writeJSON(w, http.StatusCreated)

}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := h.usecase.DeleteCustomerByID(r.Context(), usecase.DeleteCustomerRequest{
		ID: vars["id"],
	})

	if err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	writeJSON(w, http.StatusNoContent)
}

type AddFamilyMemberRequest struct {
	CustomerID string                 `json:"-"`
	List       []*AddFamilyMemberList `json:"family_list"`
}

type AddFamilyMemberList struct {
	Relation    string    `json:"relation"`
	Name        string    `json:"name"`
	Dob         string    `json:"dob"`
	DateOfBirth time.Time `json:"doobs"`
}

func (req *AddFamilyMemberRequest) Validate() error {
	if req.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}

	for _, f := range req.List {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (req *AddFamilyMemberList) Validate() error {

	if req.Relation == "" {
		return fmt.Errorf("name is required")
	}

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	v, err := time.Parse("2006-01-02", req.Dob)
	if req.Dob == "" || err != nil {
		return fmt.Errorf("dob is required")
	}
	req.DateOfBirth = v

	return nil
}

func (h *CustomerHandler) AddFamilyMember(w http.ResponseWriter, r *http.Request) {

	req := &AddFamilyMemberRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		BadRequest("invalid-request", err, w, r)
		return
	}

	vars := mux.Vars(r)
	req.CustomerID = vars["customerID"]

	// Validate the request data
	if err := req.Validate(); err != nil {
		BadRequest("validation-error", err, w, r)
		return
	}

	addFamilyMember := usecase.AddFamilyMember{
		CustomerID: req.CustomerID,
		List:       []usecase.FamilyMember{},
	}

	for _, f := range req.List {
		addFamilyMember.List = append(addFamilyMember.List, usecase.FamilyMember{
			Relation:    f.Relation,
			Name:        f.Name,
			DateOfBirth: f.DateOfBirth,
		})
	}

	if err := h.usecase.AddFamilyMember(r.Context(), addFamilyMember); err != nil {
		return
	}

	writeJSON(w, http.StatusCreated)

}

func (h *CustomerHandler) DeleteFamilyMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := h.usecase.RemoveFamilyMember(r.Context(), usecase.RemoveFamilyMember{
		CustomerID: vars["customerID"],
		FamilyID:   vars["familyMemberID"],
	})

	if err != nil {
		RespondWithSlugError(err, w, r)
		return
	}

	writeJSON(w, http.StatusNoContent)
}
