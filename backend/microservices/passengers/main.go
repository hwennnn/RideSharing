package main

import (
	"database/sql"
	"encoding/json"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	models "backend/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}

func passengers(res http.ResponseWriter, req *http.Request) {
	var results []models.Passenger
	databaseResults, err := db.Query("Select * FROM Passengers")

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		// map this type to the record in the table
		var passenger models.Passenger
		err = databaseResults.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress, &passenger.AvailableStatus)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, passenger)
		fmt.Println(passenger.PassengerID, passenger.FirstName, passenger.LastName, passenger.MobileNumber, passenger.EmailAddress)
	}

	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func passenger(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	passengerid := params["passengerid"]

	if req.Method == "GET" {
		query := fmt.Sprintf("SELECT * FROM Passengers WHERE PassengerID='%s'", passengerid)
		databaseResults, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		var isExist bool
		var passenger models.Passenger
		for databaseResults.Next() {
			err = databaseResults.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress, &passenger.AvailableStatus)
			if err != nil {
				panic(err.Error())
			}
			isExist = true
		}

		if isExist {
			json.NewEncoder(res).Encode(passenger)
		} else {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("404 - No Passenger found"))
		}

	}

	if req.Header.Get("Content-type") == "application/json" {

		// POST is for creating new driver
		if req.Method == "POST" {

			// read the string sent to the service
			var newPassenger models.Passenger
			reqBody, err := ioutil.ReadAll(req.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if !isPassengerJsonCompleted(newPassenger) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				if passengerid != newPassenger.PassengerID {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - The data in body and parameters do not match"))
					return
				}

				// check if driver exists; add only if driver does not exist
				query := fmt.Sprintf("SELECT * FROM Passengers WHERE PassengerID='%s'", passengerid)
				databaseResults, err := db.Query(query)
				if err != nil {
					panic(err.Error())
				}

				var isPassengerExist bool
				for databaseResults.Next() {
					if err != nil {
						panic(err.Error())
					}
					isPassengerExist = true
				}

				if !isPassengerExist {
					query := fmt.Sprintf("INSERT INTO Passengers VALUES ('%s', '%s', '%s', '%s', '%s', '%d')", newPassenger.PassengerID, newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.EmailAddress, 1)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}

					res.WriteHeader(http.StatusCreated)
					res.Write([]byte("201 - Passenger added: " + passengerid))
				} else {
					res.WriteHeader(http.StatusConflict)
					res.Write([]byte("409 - Duplicate passenger ID"))
				}
			} else {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply passenger information in JSON format"))
			}
		}
	}

	//---PUT is for creating or updating
	// existing course---
	if req.Method == "PUT" {
		var newPassenger models.Passenger
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newPassenger)

			if passengerid != newPassenger.PassengerID {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - The data in body and parameters do not match"))
				return
			}

			// check if driver exists; add only if driver does not exist
			query := fmt.Sprintf("SELECT * FROM Passengers WHERE PassengerID='%s'", passengerid)
			databaseResults, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}

			var isPassengerExist bool
			for databaseResults.Next() {
				if err != nil {
					panic(err.Error())
				}
				isPassengerExist = true
			}

			if !isPassengerExist {
				if !isPassengerJsonCompleted(newPassenger) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				query := fmt.Sprintf("INSERT INTO Passengers VALUES ('%s', '%s', '%s', '%s', '%s', '%d')", newPassenger.PassengerID, newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.EmailAddress, 1)

				_, err := db.Query(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusCreated)
				res.Write([]byte("201 - Passenger added: " + passengerid))
			} else {
				formattedUpdateFieldQuery := formmatedUpdatePassengerQueryField(newPassenger)

				if formattedUpdateFieldQuery == "" { // means there is no valid field can be updated
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply passenger information in JSON format"))
					return
				}

				query := fmt.Sprintf("UPDATE Passengers SET %s WHERE PassengerID='%s'", formattedUpdateFieldQuery, newPassenger.PassengerID)

				_, err := db.Query(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusAccepted)
				res.Write([]byte("202 - Passenger updated: " + passengerid))
			}

		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Please supply passenger information in JSON format"))
		}
	}
}

func formmatedUpdatePassengerQueryField(newPassenger models.Passenger) string {
	var fields []string

	if newPassenger.FirstName != "" {
		fields = append(fields, fmt.Sprintf("FirstName='%s'", newPassenger.FirstName))
	}

	if newPassenger.LastName != "" {
		fields = append(fields, fmt.Sprintf("LastName='%s'", newPassenger.LastName))
	}

	if newPassenger.MobileNumber != "" {
		fields = append(fields, fmt.Sprintf("MobileNumber='%s'", newPassenger.MobileNumber))
	}

	if newPassenger.EmailAddress != "" {
		fields = append(fields, fmt.Sprintf("EmailAddress='%s'", newPassenger.EmailAddress))
	}

	if newPassenger.AvailableStatus != 0 {
		fields = append(fields, fmt.Sprintf("AvailableStatus='%d'", newPassenger.AvailableStatus))
	}

	return strings.Join(fields, ", ")
}

func isPassengerJsonCompleted(passenger models.Passenger) bool {
	passengerID := strings.TrimSpace(passenger.PassengerID)
	firstName := strings.TrimSpace(passenger.FirstName)
	lastName := strings.TrimSpace(passenger.LastName)
	mobileNumber := strings.TrimSpace(passenger.MobileNumber)
	emailAddress := strings.TrimSpace(passenger.EmailAddress)

	fmt.Println(passengerID, firstName, lastName, mobileNumber, emailAddress)
	return passengerID != "" && firstName != "" && lastName != "" && mobileNumber != "" && emailAddress != ""
}

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/api/v1/passengers/", passengers).Methods("GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", passenger).Methods("GET", "PUT", "POST")

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Passengers database server -- Listening at port 8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
