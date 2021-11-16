package main

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Driver struct {
	DriverID             int64  `json:"driver_id"`
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
	PassengerID  int64  `json:"passenger_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	EmailAddress string `json:"email_address"`
}

type Trip struct {
	TripID            int64     `json:"trip_id"`
	PassengerID       int64     `json:"passenger_id"`
	DriverID          int64     `json:"driver_id"`
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
	results := map[int64]Driver{}
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
		results[driver.DriverID] = driver
		fmt.Println(driver.DriverID, driver.FirstName, driver.LastName, driver.MobileNumber, driver.EmailAddress, driver.IdentificationNumber, driver.CarLicenseNumber, driver.AvailableStatus)
	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func passengers(res http.ResponseWriter, req *http.Request) {
	results := map[int64]Passenger{}
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
		results[passenger.PassengerID] = passenger
		fmt.Println(passenger.PassengerID, passenger.FirstName, passenger.LastName, passenger.MobileNumber, passenger.EmailAddress)
	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func trips(res http.ResponseWriter, req *http.Request) {
	results := map[int64]Trip{}

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
		results[trip.TripID] = trip

	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

// func course(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	fmt.Println(params["courseid"])

// 	// params := mux.Vars(r)

// 	/*
// 		fmt.Fprintf(w, "Detail for course "+params["courseid"])
// 		fmt.Fprintf(w, "\n")
// 		fmt.Fprintf(w, r.Method)
// 	*/

// 	if r.Method == "GET" {
// 		if _, ok := courses[params["courseid"]]; ok {
// 			json.NewEncoder(w).Encode(
// 				courses[params["courseid"]])
// 		} else {
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte("404 - No course found"))
// 		}
// 	}

// 	if r.Method == "DELETE" {
// 		if _, ok := courses[params["courseid"]]; ok {
// 			delete(courses, params["courseid"])
// 			w.WriteHeader(http.StatusAccepted)
// 			w.Write([]byte("202 - Course deleted: " +
// 				params["courseid"]))
// 		} else {
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte("404 - No course found"))
// 		}
// 	}

// 	if r.Header.Get("Content-type") == "application/json" {

// 		// POST is for creating new course
// 		if r.Method == "POST" {

// 			// read the string sent to the service
// 			var newCourse courseInfo
// 			reqBody, err := ioutil.ReadAll(r.Body)

// 			if err == nil {
// 				// convert JSON to object
// 				json.Unmarshal(reqBody, &newCourse)

// 				if newCourse.Title == "" {
// 					w.WriteHeader(
// 						http.StatusUnprocessableEntity)
// 					w.Write([]byte(
// 						"422 - Please supply course " +
// 							"information " + "in JSON format"))
// 					return
// 				}

// 				// check if course exists; add only if
// 				// course does not exist
// 				if _, ok := courses[params["courseid"]]; !ok {
// 					courses[params["courseid"]] = newCourse
// 					w.WriteHeader(http.StatusCreated)
// 					w.Write([]byte("201 - Course added: " +
// 						params["courseid"]))
// 				} else {
// 					w.WriteHeader(http.StatusConflict)
// 					w.Write([]byte(
// 						"409 - Duplicate course ID"))
// 				}
// 			} else {
// 				w.WriteHeader(
// 					http.StatusUnprocessableEntity)
// 				w.Write([]byte("422 - Please supply course information " +
// 					"in JSON format"))
// 			}
// 		}

// 		//---PUT is for creating or updating
// 		// existing course---
// 		if r.Method == "PUT" {
// 			var newCourse courseInfo
// 			reqBody, err := ioutil.ReadAll(r.Body)

// 			if err == nil {
// 				json.Unmarshal(reqBody, &newCourse)

// 				if newCourse.Title == "" {
// 					w.WriteHeader(
// 						http.StatusUnprocessableEntity)
// 					w.Write([]byte(
// 						"422 - Please supply course " +
// 							" information " +
// 							"in JSON format"))
// 					return
// 				}

// 				// check if course exists; add only if
// 				// course does not exist
// 				if _, ok := courses[params["courseid"]]; !ok {
// 					courses[params["courseid"]] =
// 						newCourse
// 					w.WriteHeader(http.StatusCreated)
// 					w.Write([]byte("201 - Course added: " +
// 						params["courseid"]))
// 				} else {
// 					// update course
// 					courses[params["courseid"]] = newCourse
// 					w.WriteHeader(http.StatusAccepted)
// 					w.Write([]byte("202 - Course updated: " +
// 						params["courseid"]))
// 				}
// 			} else {
// 				w.WriteHeader(
// 					http.StatusUnprocessableEntity)
// 				w.Write([]byte("422 - Please supply " +
// 					"course information " +
// 					"in JSON format"))
// 			}
// 		}
// 	}
// }

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/drivers/", drivers)
	router.HandleFunc("/api/v1/passengers/", passengers)
	router.HandleFunc("/api/v1/trips/", trips)
	// router.HandleFunc("/api/v1/courses", allcourses)
	// router.HandleFunc("/api/v1/courses/{courseid}", course).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
