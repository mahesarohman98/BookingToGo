package http

import "github.com/gorilla/mux"

func NewRouter(customerHandler *CustomerHandler) *mux.Router {
	router := mux.NewRouter()

	// CRUD operations for Customer
	router.HandleFunc("/customers", customerHandler.CreateCustomer).Methods("POST")
	router.HandleFunc("/customers", customerHandler.GetListCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", customerHandler.GetCustomerByID).Methods("GET")
	router.HandleFunc("/customers/{id}", customerHandler.UpdateCustomer).Methods("PATCH")
	router.HandleFunc("/customers/{id}", customerHandler.DeleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{customerID}/family", customerHandler.AddFamilyMember).Methods("POST")
	router.HandleFunc("/customers/{customerID}/family/{familyMemberID}", customerHandler.DeleteFamilyMember).Methods("DELETE")

	router.HandleFunc("/countries", customerHandler.GetCountries).Methods("Get")
	return router
}
