package main

import (
	"booking_to_go/internal/interface/config"
	httpInterface "booking_to_go/internal/interface/http"
	"booking_to_go/internal/interface/repositories"
	"booking_to_go/internal/usecase"
	"fmt"
	"net/http"
	"os"
)

func main() {

	conn, err := config.NewPostgreSQLConnection()
	if err != nil {
		panic(fmt.Sprint("Error connecting db:", err))
	}
	defer conn.Close()

	repository := repositories.NewCustomerPGRepository(conn)

	usecase := usecase.NewCustomerUsecase(repository)

	handler := httpInterface.NewCustomerHandler(usecase)

	// Set up Gorilla Mux router
	router := httpInterface.NewRouter(handler)

	// Start HTTP server
	fmt.Println("Starting server on :", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		panic(fmt.Sprint("Error starting server:", err))
	}

}
