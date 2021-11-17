package main

import (
	"database/sql"
	"encoding/json"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Driver struct {
	DriverID             string `json:"driver_id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	MobileNumber         string `json:"mobile_number"`
	EmailAddress         string `json:"email_address"`
	IdentificationNumber string `json:"identification_number"`
	CarLicenseNumber     string `json:"car_license_number"`
	AvailableStatus      int    `json:"available_status"`
	// 0 => Offline
	// 1 => Online and available
	// 2 => Online and during the trip
}

type Passenger struct {
	PassengerID  string `json:"passenger_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	EmailAddress string `json:"email_address"`
}

type Trip struct {
	TripID            string    `json:"trip_id"`
	PassengerID       string    `json:"passenger_id"`
	DriverID          string    `json:"driver_id"`
	PickupPostalCode  string    `json:"pickup_postal_code"`
	DropoffPostalCode string    `json:"dropoff_postal_code"`
	TripProgress      int       `json:"trip_progress"`
	Passenger         Passenger `json:"passenger"`
	Driver            Driver    `json:"driver"`

	// 0 -> Created by passenger, but no driver is assigned yet
	// 1 -> The trip is ongoing
	// 2 -> The trip has ended
}

var db *sql.DB

func home(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to the REST API!")
}

func drivers(res http.ResponseWriter, req *http.Request) {
	var results []Driver
	databaseResults, err := db.Query("Select * FROM Drivers")

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		// map this type to the record in the table
		var driver Driver
		err = databaseResults.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.IdentificationNumber, &driver.CarLicenseNumber, &driver.AvailableStatus)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, driver)
		fmt.Println(driver.DriverID, driver.FirstName, driver.LastName, driver.MobileNumber, driver.EmailAddress, driver.IdentificationNumber, driver.CarLicenseNumber, driver.AvailableStatus)
	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func passengers(res http.ResponseWriter, req *http.Request) {
	var results []Passenger
	databaseResults, err := db.Query("Select * FROM Passengers")

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		// map this type to the record in the table
		var passenger Passenger
		err = databaseResults.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, passenger)
		fmt.Println(passenger.PassengerID, passenger.FirstName, passenger.LastName, passenger.MobileNumber, passenger.EmailAddress)
	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func trips(res http.ResponseWriter, req *http.Request) {
	var results []Trip

	databaseResults, err := db.Query("SELECT * FROM Trips t INNER JOIN Drivers d ON t.DriverID = d.DriverID INNER JOIN Passengers p ON t.PassengerID = p.PassengerID")
	// [TripID PassengerID DriverID PickupPostalCode DropoffPostalCode TripProgress DriverID FirstName LastName MobileNumber EmailAddress IdentificationNumber CarLicenseNumber AvailableStatus PassengerID FirstName LastName MobileNumber EmailAddress]
	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		var trip Trip
		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &trip.DriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.Driver.DriverID, &trip.Driver.FirstName, &trip.Driver.LastName, &trip.Driver.MobileNumber, &trip.Driver.EmailAddress, &trip.Driver.IdentificationNumber, &trip.Driver.CarLicenseNumber, &trip.Driver.AvailableStatus, &trip.Passenger.PassengerID, &trip.Passenger.FirstName, &trip.Passenger.LastName, &trip.Passenger.MobileNumber, &trip.Passenger.EmailAddress)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, trip)

	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func driver(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	driverid := params["driverid"]

	if req.Method == "GET" {
		query := fmt.Sprintf("SELECT * FROM Drivers WHERE DriverID='%s'", driverid)
		databaseResults, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		var isExist bool
		var driver Driver
		for databaseResults.Next() {
			err = databaseResults.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.IdentificationNumber, &driver.CarLicenseNumber, &driver.AvailableStatus)
			if err != nil {
				panic(err.Error())
			}
			isExist = true
		}

		if isExist {
			json.NewEncoder(res).Encode(driver)
		} else {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("404 - No driver found"))
		}

	}

	if req.Header.Get("Content-type") == "application/json" {

		// POST is for creating new driver
		if req.Method == "POST" {

			// read the string sent to the service
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(req.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)

				if !isDriverJsonCompleted(newDriver) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				if driverid != newDriver.DriverID {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - The data in body and parameters do not match"))
					return
				}

				// check if driver exists; add only if driver does not exist
				query := fmt.Sprintf("SELECT * FROM Drivers WHERE DriverID='%s'", driverid)
				databaseResults, err := db.Query(query)
				if err != nil {
					panic(err.Error())
				}

				var isDriverExist bool
				for databaseResults.Next() {
					if err != nil {
						panic(err.Error())
					}
					isDriverExist = true
				}

				if !isDriverExist {
					query := fmt.Sprintf("INSERT INTO Drivers VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %d)", newDriver.DriverID, newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.EmailAddress, newDriver.IdentificationNumber, newDriver.CarLicenseNumber, 0)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}

					res.WriteHeader(http.StatusCreated)
					res.Write([]byte("201 - Driver added: " + driverid))
				} else {
					res.WriteHeader(http.StatusConflict)
					res.Write([]byte("409 - Duplicate driver ID"))
				}
			} else {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply driver information in JSON format"))
			}
		}
	}

	// 	//---PUT is for creating or updating
	// 	// existing course---
	// 	if r.Method == "PUT" {
	// 		var newDriver courseInfo
	// 		reqBody, err := ioutil.ReadAll(r.Body)

	// 		if err == nil {
	// 			json.Unmarshal(reqBody, &newDriver)

	// 			if newDriver.Title == "" {
	// 				w.WriteHeader(
	// 					http.StatusUnprocessableEntity)
	// 				w.Write([]byte(
	// 					"422 - Please supply course " +
	// 						" information " +
	// 						"in JSON format"))
	// 				return
	// 			}

	// 			// check if course exists; add only if
	// 			// course does not exist
	// 			if _, ok := courses[params["courseid"]]; !ok {
	// 				courses[params["courseid"]] =
	// 					newDriver
	// 				w.WriteHeader(http.StatusCreated)
	// 				w.Write([]byte("201 - Course added: " +
	// 					params["courseid"]))
	// 			} else {
	// 				// update course
	// 				courses[params["courseid"]] = newDriver
	// 				w.WriteHeader(http.StatusAccepted)
	// 				w.Write([]byte("202 - Course updated: " +
	// 					params["courseid"]))
	// 			}
	// 		} else {
	// 			w.WriteHeader(
	// 				http.StatusUnprocessableEntity)
	// 			w.Write([]byte("422 - Please supply " +
	// 				"course information " +
	// 				"in JSON format"))
	// 		}
	// 	}
	// }
}

func isDriverJsonCompleted(driver Driver) bool {
	driverID := strings.TrimSpace(driver.DriverID)
	firstName := strings.TrimSpace(driver.FirstName)
	lastName := strings.TrimSpace(driver.LastName)
	mobileNumber := strings.TrimSpace(driver.MobileNumber)
	emailAddress := strings.TrimSpace(driver.EmailAddress)
	identificationNumber := strings.TrimSpace(driver.IdentificationNumber)
	carLicenseNumber := strings.TrimSpace(driver.CarLicenseNumber)
	fmt.Println(driverID, firstName, lastName, mobileNumber, emailAddress, identificationNumber, carLicenseNumber)
	return driverID != "" && firstName != "" && lastName != "" && mobileNumber != "" && emailAddress != "" && identificationNumber != "" && carLicenseNumber != ""
}

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/drivers/", drivers)
	router.HandleFunc("/api/v1/driver/{driverid}", driver).Methods("GET", "PUT", "POST")
	router.HandleFunc("/api/v1/passengers/", passengers)
	router.HandleFunc("/api/v1/trips/", trips)
	// router.HandleFunc("/api/v1/courses", allcourses)
	// router.HandleFunc("/api/v1/courses/{courseid}", course).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
