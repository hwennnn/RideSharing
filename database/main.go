package main

import (
	"database/sql"
	"encoding/json"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	// 0 -> used by golang to indicate whether the integer variable has been initialised or not
	// 1 => Online and available
	// 2 => Online but during the trip
}

type Passenger struct {
	PassengerID     string `json:"passenger_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MobileNumber    string `json:"mobile_number"`
	EmailAddress    string `json:"email_address"`
	AvailableStatus int    `json:"available_status"`
	// 0 -> used by golang to indicate whether the integer variable has been initialised or not
	// 1 => Online and available
	// 2 => Online but during the trip
}

type Trip struct {
	TripID            string    `json:"trip_id"`
	PassengerID       string    `json:"passenger_id"`
	DriverID          string    `json:"driver_id"`
	PickupPostalCode  string    `json:"pickup_postal_code"`
	DropoffPostalCode string    `json:"dropoff_postal_code"`
	TripProgress      int       `json:"trip_progress"`
	CreatedTime       int64     `json:"created_time"`
	CompletedTime     int64     `json:"completed_time"`
	Passenger         Passenger `json:"passenger"`
	Driver            Driver    `json:"driver"`
	// 0 -> used by golang to indicate whether the integer variable has been initialised or not
	// 1 -> Created by passenger, but no driver is assigned yet
	// 2 -> The trip is ongoing
	// 3 -> The trip has ended
}

var db *sql.DB

func currentMs() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}

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

func availableDrivers(res http.ResponseWriter, req *http.Request) {
	var results []Driver
	databaseResults, err := db.Query("Select * FROM Drivers WHERE AvailableStatus = 1")

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

	if len(results) == 0 {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("404 - No available drivers are found"))
	} else {
		// returns the first available driver in JSON
		json.NewEncoder(res).Encode(results[0])
	}
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
					query := fmt.Sprintf("INSERT INTO Drivers VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %d)", newDriver.DriverID, newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.EmailAddress, newDriver.IdentificationNumber, newDriver.CarLicenseNumber, 1)

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

	//---PUT is for creating or updating
	// existing course---
	if req.Method == "PUT" {
		var newDriver Driver
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newDriver)

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
				if !isDriverJsonCompleted(newDriver) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				query := fmt.Sprintf("INSERT INTO Drivers VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %d)", newDriver.DriverID, newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.EmailAddress, newDriver.IdentificationNumber, newDriver.CarLicenseNumber, 1)

				_, err := db.Query(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusCreated)
				res.Write([]byte("201 - Driver added: " + driverid))
			} else {
				formattedUpdateFieldQuery := formmatedUpdateDriverQueryField(newDriver)

				if formattedUpdateFieldQuery == "" { // means there is no valid field can be updated
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply driver information in JSON format"))
					return
				}

				query := fmt.Sprintf("UPDATE Drivers SET %s WHERE DriverID='%s'", formattedUpdateFieldQuery, newDriver.DriverID)

				_, err := db.Query(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusAccepted)
				res.Write([]byte("202 - Driver updated: " + driverid))
			}

		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Please supply driver information in JSON format"))
		}
	}
}

func formmatedUpdateDriverQueryField(newDriver Driver) string {
	var fields []string

	if newDriver.FirstName != "" {
		fields = append(fields, fmt.Sprintf("FirstName='%s'", newDriver.FirstName))
	}

	if newDriver.LastName != "" {
		fields = append(fields, fmt.Sprintf("LastName='%s'", newDriver.LastName))
	}

	if newDriver.MobileNumber != "" {
		fields = append(fields, fmt.Sprintf("MobileNumber='%s'", newDriver.MobileNumber))
	}

	if newDriver.EmailAddress != "" {
		fields = append(fields, fmt.Sprintf("EmailAddress='%s'", newDriver.EmailAddress))
	}

	if newDriver.CarLicenseNumber != "" {
		fields = append(fields, fmt.Sprintf("CarLicenseNumber='%s'", newDriver.CarLicenseNumber))
	}

	if newDriver.AvailableStatus != 0 {
		fields = append(fields, fmt.Sprintf("AvailableStatus='%d'", newDriver.AvailableStatus))
	}

	return strings.Join(fields, ", ")
}

func isDriverJsonCompleted(driver Driver) bool {
	driverID := strings.TrimSpace(driver.DriverID)
	firstName := strings.TrimSpace(driver.FirstName)
	lastName := strings.TrimSpace(driver.LastName)
	mobileNumber := strings.TrimSpace(driver.MobileNumber)
	emailAddress := strings.TrimSpace(driver.EmailAddress)
	identificationNumber := strings.TrimSpace(driver.IdentificationNumber)
	carLicenseNumber := strings.TrimSpace(driver.CarLicenseNumber)
	// fmt.Println(driverID, firstName, lastName, mobileNumber, emailAddress, identificationNumber, carLicenseNumber)
	return driverID != "" && firstName != "" && lastName != "" && mobileNumber != "" && emailAddress != "" && identificationNumber != "" && carLicenseNumber != ""
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
		var passenger Passenger
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
			var newPassenger Passenger
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
		var newPassenger Passenger
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

func formmatedUpdatePassengerQueryField(newPassenger Passenger) string {
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

func isPassengerJsonCompleted(passenger Passenger) bool {
	passengerID := strings.TrimSpace(passenger.PassengerID)
	firstName := strings.TrimSpace(passenger.FirstName)
	lastName := strings.TrimSpace(passenger.LastName)
	mobileNumber := strings.TrimSpace(passenger.MobileNumber)
	emailAddress := strings.TrimSpace(passenger.EmailAddress)

	fmt.Println(passengerID, firstName, lastName, mobileNumber, emailAddress)
	return passengerID != "" && firstName != "" && lastName != "" && mobileNumber != "" && emailAddress != ""
}

type TripsRequestBody struct {
	PassengerID string `json:"passenger_id"`
	DriverID    string `json:"driver_id"`
}

func trips(res http.ResponseWriter, req *http.Request) {
	var results []Trip

	params := req.URL.Query()

	// convert JSON to object
	formmatedFieldQuery := formmatedTripQueryField(params["driver_id"], params["passenger_id"])

	query := fmt.Sprintf("SELECT * FROM Trips t INNER JOIN Drivers d ON t.DriverID = d.DriverID INNER JOIN Passengers p ON t.PassengerID = p.PassengerID WHERE %s", formmatedFieldQuery)
	fmt.Println(query)
	databaseResults, err := db.Query(query)

	// [TripID PassengerID DriverID PickupPostalCode DropoffPostalCode TripProgress DriverID FirstName LastName MobileNumber EmailAddress IdentificationNumber CarLicenseNumber AvailableStatus PassengerID FirstName LastName MobileNumber EmailAddress]
	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		var trip Trip
		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &trip.DriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime, &trip.Driver.DriverID, &trip.Driver.FirstName, &trip.Driver.LastName, &trip.Driver.MobileNumber, &trip.Driver.EmailAddress, &trip.Driver.IdentificationNumber, &trip.Driver.CarLicenseNumber, &trip.Driver.AvailableStatus, &trip.Passenger.PassengerID, &trip.Passenger.FirstName, &trip.Passenger.LastName, &trip.Passenger.MobileNumber, &trip.Passenger.EmailAddress)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, trip)

	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func formmatedTripQueryField(driverID []string, passengerID []string) string {
	var results string
	fmt.Println(driverID, passengerID)

	if len(driverID) > 0 && driverID[0] != "" {
		results += fmt.Sprintf("t.DriverID = '%s'", driverID[0])
	}

	if len(passengerID) > 0 && passengerID[0] != "" {
		if results != "" {
			results += " AND "
		}

		results += fmt.Sprintf("t.PassengerID = '%s'", passengerID[0])
	}

	if results == "" {
		return "t.TripProgress=3"
	}

	return results + " AND t.TripProgress=3" // filter only completed trip
}

func trip(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	tripid := params["tripid"]

	if req.Method == "GET" {
		query := fmt.Sprintf("SELECT * FROM Trips t INNER JOIN Drivers d ON t.DriverID = d.DriverID INNER JOIN Passengers p ON t.PassengerID = p.PassengerID WHERE TripID='%s'", tripid)
		databaseResults, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		var isExist bool
		var trip Trip
		for databaseResults.Next() {
			err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &trip.DriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime, &trip.Driver.DriverID, &trip.Driver.FirstName, &trip.Driver.LastName, &trip.Driver.MobileNumber, &trip.Driver.EmailAddress, &trip.Driver.IdentificationNumber, &trip.Driver.CarLicenseNumber, &trip.Driver.AvailableStatus, &trip.Passenger.PassengerID, &trip.Passenger.FirstName, &trip.Passenger.LastName, &trip.Passenger.MobileNumber, &trip.Passenger.EmailAddress)
			if err != nil {
				panic(err.Error())
			}
			isExist = true
		}

		if isExist {
			json.NewEncoder(res).Encode(trip)
		} else {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("404 - No Trip found"))
		}

	}

	if req.Header.Get("Content-type") == "application/json" {

		// POST is for creating new driver
		if req.Method == "POST" {

			// read the string sent to the service
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(req.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newTrip)

				if !isTripJsonCompleted(newTrip) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply trip information in JSON format"))
					return
				}

				if tripid != newTrip.TripID {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - The data in body and parameters do not match"))
					return
				}

				// check if driver exists; add only if driver does not exist
				query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID='%s'", tripid)
				databaseResults, err := db.Query(query)
				if err != nil {
					panic(err.Error())
				}

				var isTripExist bool
				for databaseResults.Next() {
					if err != nil {
						panic(err.Error())
					}
					isTripExist = true
				}

				if !isTripExist {
					query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', NULL, '%s', '%s', '%d', '%d', NULL)", newTrip.TripID, newTrip.PassengerID, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, 1, currentMs())

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}

					res.WriteHeader(http.StatusCreated)
					res.Write([]byte("201 - Trip added: " + tripid))
				} else {
					res.WriteHeader(http.StatusConflict)
					res.Write([]byte("409 - Duplicate trip ID"))
				}
			} else {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply trip information in JSON format"))
			}
		}
	}

	//---PUT is for creating or updating
	// existing course---
	if req.Method == "PUT" {
		var newTrip Trip
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newTrip)

			if tripid != newTrip.TripID {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - The data in body and parameters do not match"))
				return
			}

			// check if trip exists; add only if trip does not exist
			query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID='%s'", tripid)
			databaseResults, err := db.Query(query)
			if err != nil {
				panic(err.Error())
			}

			var tripFromDatabase Trip
			var isTripExist bool
			for databaseResults.Next() {
				if err != nil {
					panic(err.Error())
				}
				err = databaseResults.Scan(&tripFromDatabase.TripID, &tripFromDatabase.PassengerID, &tripFromDatabase.DriverID, &tripFromDatabase.PickupPostalCode, &tripFromDatabase.DropoffPostalCode, &tripFromDatabase.TripProgress, &tripFromDatabase.CreatedTime, &tripFromDatabase.CompletedTime, &tripFromDatabase.Driver.DriverID, &tripFromDatabase.Driver.FirstName, &tripFromDatabase.Driver.LastName, &tripFromDatabase.Driver.MobileNumber, &tripFromDatabase.Driver.EmailAddress, &tripFromDatabase.Driver.IdentificationNumber, &tripFromDatabase.Driver.CarLicenseNumber, &tripFromDatabase.Driver.AvailableStatus, &tripFromDatabase.Passenger.PassengerID, &tripFromDatabase.Passenger.FirstName, &tripFromDatabase.Passenger.LastName, &tripFromDatabase.Passenger.MobileNumber, &tripFromDatabase.Passenger.EmailAddress)
				isTripExist = true
			}

			if !isTripExist {
				if !isTripJsonCompleted(newTrip) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply trip information in JSON format"))
					return
				}

				query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', NULL)", newTrip.TripID, newTrip.PassengerID, newTrip.DriverID, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, 1, currentMs())

				_, err := db.Query(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusCreated)
				res.Write([]byte("201 - Trip added: " + tripid))
			} else {
				formattedUpdateFieldQuery := formattedUpdateTripQueryField(newTrip)

				if formattedUpdateFieldQuery == "" { // means there is no valid field can be updated
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply trip information in JSON format"))
					return
				}

				if tripFromDatabase.TripProgress != 3 && newTrip.TripProgress == 3 {
					updateTripCompletedTime(tripid)
				}

				query := fmt.Sprintf("UPDATE Trips SET %s WHERE TripID='%s'", formattedUpdateFieldQuery, newTrip.TripID)

				_, err := db.Query(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusAccepted)
				res.Write([]byte("202 - Trip updated: " + tripid))
			}

		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Please supply trip information in JSON format"))
		}
	}
}

func updateTripCompletedTime(tripid string) {
	query := fmt.Sprintf("UPDATE Trips SET CompletedTime=%d WHERE TripID=%s", currentMs(), tripid)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func formattedUpdateTripQueryField(trip Trip) string {
	var results string

	if trip.DriverID != "" {
		results += fmt.Sprintf("DriverID='%s'", trip.DriverID)
	}

	if trip.TripProgress != 0 {
		results += fmt.Sprintf("TripProgress='%d'", trip.TripProgress)
	}

	return results
}

func isTripJsonCompleted(trip Trip) bool {
	tripID := strings.TrimSpace(trip.TripID)
	passengerID := strings.TrimSpace(trip.PassengerID)
	driverID := strings.TrimSpace(trip.DriverID)
	pickupPostalCode := strings.TrimSpace(trip.PickupPostalCode)
	dropoffPostalCode := strings.TrimSpace(trip.DropoffPostalCode)

	return tripID != "" && passengerID != "" && driverID != "" && pickupPostalCode != "" && dropoffPostalCode != ""
}

// func setupCorsResponse(res *http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Access-Control-Allow-Origin", "*")
// 	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
// }

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// TODO: Authenticate requests by token
		// token := req.Header.Get("X-Session-Token")

		// if user, found := amw.tokenUsers[token]; found {
		// 	// We found the token in our map
		// 	log.Printf("Authenticated user %s\n", user)
		// 	next.ServeHTTP(w, r)
		// } else {
		// 	http.Error(w, "Forbidden", http.StatusForbidden)
		// }
		res.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/api/v1/", home).Methods("GET")

	router.HandleFunc("/api/v1/drivers/", drivers).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{driverid}", driver).Methods("GET", "PUT", "POST")
	router.HandleFunc("/api/v1/available_drivers/", availableDrivers).Methods("GET")

	router.HandleFunc("/api/v1/passengers/", passengers).Methods("GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", passenger).Methods("GET", "PUT", "POST")

	router.HandleFunc("/api/v1/trips", trips).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods("GET", "PUT", "POST")

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}
