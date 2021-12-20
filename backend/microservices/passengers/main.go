package main

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	models "passengers/models"
	utils "passengers/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// global database handler object
var db *sql.DB

// this middleware will set the returned content type as application/json
// this helps reduce code redudancy, which originally has to be added in each response writer
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}

// This method is used to retrieve passengers from MySQL,
// and return the result in array of passenger json object
func getPassengers(res http.ResponseWriter, req *http.Request) {
	var results []models.Passenger

	databaseResults, err := db.Query("Select * FROM Passengers")

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		// Map the passenger object to the record in the table

		var passenger models.Passenger
		err = databaseResults.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress, &passenger.AvailableStatus)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, passenger)
	}

	// Returns all the passengers in JSON
	json.NewEncoder(res).Encode(results)
}

// This method is used to retrieve a passenger from MySQL by specific passengerID,
// and return the result in json otherwise return 404 code
func getPassenger(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	passengerid := params["passengerid"]

	isPassengerExist, passenger := getPassengerHelper(passengerid)

	if isPassengerExist {
		json.NewEncoder(res).Encode(passenger)
	} else {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("404 - No Passenger found"))
	}
}

// This helper method helps to query the passenger from the database,
// and return (boolean, passenger) tuple object
func getPassengerHelper(passengerID string) (bool, models.Passenger) {
	query := fmt.Sprintf("SELECT * FROM Passengers WHERE PassengerID='%s'", passengerID)
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

	return isExist, passenger
}

// This method is used to create a passenger in MySQL by specific passengerID,
// Case 1: If the compulsory passenger information is not provided, it will return message which says the information is not correctly supplied
// Case 2: It will fail and return conflict status code if a passenger with same passengerID is already found in the database
// Case 3: Otherwise, it will return success message with status created code
func postPassenger(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	passengerid := params["passengerid"]

	// read the body string sent to the service
	var newPassenger models.Passenger
	reqBody, err := ioutil.ReadAll(req.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &newPassenger)

		if !utils.IsPassengerJsonCompleted(newPassenger) {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Please supply passenger information in JSON format"))
			return
		}

		if passengerid != newPassenger.PassengerID {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - The data in body and parameters do not match"))
			return
		}

		// check if passenger exists; add only if passenger does not exist
		isPassengerExist, _ := getPassengerHelper(passengerid)

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

// This method is used for either creating or updating the passenger depends whether the passengerID exists
// Case 1: If passengerID exists, update the passenger using the information retrieved from request body
// Case 2: If passengerID does not exist, create the passenger using the information retrieved from request body
func putPassenger(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	passengerid := params["passengerid"]

	var newPassenger models.Passenger
	reqBody, err := ioutil.ReadAll(req.Body)

	if err == nil {
		json.Unmarshal(reqBody, &newPassenger)

		if passengerid != newPassenger.PassengerID {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - The data in body and parameters do not match"))
			return
		}

		// check if passenger exists; add only if passenger does not exist, else update
		isPassengerExist, _ := getPassengerHelper(passengerid)

		if !isPassengerExist {
			if !utils.IsPassengerJsonCompleted(newPassenger) {
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
			formattedUpdateFieldQuery := utils.FormmatedUpdatePassengerQueryField(newPassenger)

			// means there is no valid field can be updated
			if formattedUpdateFieldQuery == "" {
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

func main() {

	// Use mysql as passengerName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(db:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/api/v1/passengers/", getPassengers).Methods("GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", getPassenger).Methods("GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", postPassenger).Methods("POST")
	router.HandleFunc("/api/v1/passengers/{passengerid}", putPassenger).Methods("PUT")

	// enable cross-origin resource sharing (cors) for all requests
	handler := cors.AllowAll().Handler(router)

	fmt.Println("Passengers database server -- Listening at port 8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
