package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	TripID            string `json:"trip_id"`
	PassengerID       string `json:"passenger_id"`
	sqlDriverID       sql.NullString
	DriverID          string    `json:"driver_id"`
	PickupPostalCode  string    `json:"pickup_postal_code"`
	DropoffPostalCode string    `json:"dropoff_postal_code"`
	TripProgress      int       `json:"trip_progress"`
	CreatedTime       int64     `json:"created_time"`
	CompletedTime     int64     `json:"completed_time"`
	Passenger         Passenger `json:"passenger"`
	Driver            Driver    `json:"driver"`
	// 0 -> used by golang to indicate whether the integer variable has been initialised or not
	// 1 -> Created by passenger, but no driver is found to be assgined
	// 2 -> A driver was already assigned for the trip, but the driver has not inititated the trip yet
	// 3 -> The trip is ongoing (Driver has initiated the trip)
	// 4 -> The trip has ended (Driver has ended the trip)
}

var db *sql.DB

const authenticationToken = "1467a2a8-fff7-45b5-986d-679382d0707a"

func currentMs() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}

type TripsRequestBody struct {
	PassengerID string `json:"passenger_id"`
	DriverID    string `json:"driver_id"`
}

func fetchDriver(driverID string) Driver {
	var result Driver

	url := fmt.Sprintf("http://localhost:5000/api/v1/drivers/%s", driverID)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + authenticationToken

	// Create a new request using http
	req, _ := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result
}

func fetchPassenger(passengerID string) Passenger {
	var result Passenger

	url := fmt.Sprintf("http://localhost:5000/api/v1/passengers/%s", passengerID)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + authenticationToken

	// Create a new request using http
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result
}

func trips(res http.ResponseWriter, req *http.Request) {
	var results []Trip

	params := req.URL.Query()

	formmatedFieldQuery := formmatedTripQueryField(params["driver_id"], params["passenger_id"], params["trip_progress"])
	query := fmt.Sprintf("SELECT * FROM Trips %s ORDER BY CompletedTime DESC", formmatedFieldQuery)
	fmt.Println(query)
	databaseResults, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		var trip Trip
		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &trip.sqlDriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime)
		if trip.sqlDriverID.Valid {
			trip.DriverID = trip.sqlDriverID.String
		}

		if trip.DriverID != "" {
			trip.Driver = fetchDriver(trip.DriverID)
		}

		if trip.PassengerID != "" {
			trip.Passenger = fetchPassenger(trip.PassengerID)
		}

		if err != nil {
			panic(err.Error())
		}
		results = append(results, trip)

	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func formmatedTripQueryField(driverID []string, passengerID []string, tripProgress []string) string {
	var results string

	if len(driverID) > 0 && driverID[0] != "" {
		results += fmt.Sprintf("DriverID = '%s'", driverID[0])
	}

	if len(passengerID) > 0 && passengerID[0] != "" {
		if results != "" {
			results += " AND "
		}

		results += fmt.Sprintf("PassengerID = '%s'", passengerID[0])
	}

	if len(tripProgress) > 0 && tripProgress[0] != "" {
		if results != "" {
			results += " AND "
		}
		parsedTripProgress, _ := strconv.ParseInt(tripProgress[0], 10, 64)

		results += fmt.Sprintf("TripProgress = %d", parsedTripProgress)
	}

	if results == "" {
		return ""
	}

	return "WHERE " + results
}

func trip(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	tripid := params["tripid"]

	if req.Method == "GET" {
		query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID='%s'", tripid)
		databaseResults, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		var isExist bool
		var trip Trip
		for databaseResults.Next() {
			err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &trip.sqlDriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime)

			if trip.sqlDriverID.Valid {
				trip.DriverID = trip.sqlDriverID.String
			}

			if trip.DriverID != "" {
				trip.Driver = fetchDriver(trip.DriverID)
			}

			if trip.PassengerID != "" {
				trip.Passenger = fetchPassenger(trip.PassengerID)
			}

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

				createTrip(newTrip, res, req)

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
				err = databaseResults.Scan(&tripFromDatabase.TripID, &tripFromDatabase.PassengerID, &tripFromDatabase.DriverID, &tripFromDatabase.PickupPostalCode, &tripFromDatabase.DropoffPostalCode, &tripFromDatabase.TripProgress, &tripFromDatabase.CreatedTime, &tripFromDatabase.CompletedTime)
				isTripExist = true
			}

			if !isTripExist {
				if !isTripJsonCompleted(newTrip) {
					res.WriteHeader(http.StatusUnprocessableEntity)
					res.Write([]byte("422 - Please supply trip information in JSON format"))
					return
				}

				query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', NULL, '%s', '%s', '%d', '%d', '%d')", newTrip.TripID, newTrip.PassengerID, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, 1, currentMs(), 0)

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
					// update driver available status to 3 as driver has initiated the trip and currently during the trip
					updateDriverAvailableStatus(3, newTrip.DriverID)
				} else if tripFromDatabase.TripProgress != 4 && newTrip.TripProgress == 4 {
					// update driver and passenger available status back to 1 once the trip is completed
					updateDriverAvailableStatus(1, tripFromDatabase.DriverID)
					updatePassengerAvailableStatus(1, tripFromDatabase.PassengerID)
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

// check if trip exists by tripid
func isTripExist(tripid string) bool {
	exist := false
	query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID='%s'", tripid)
	databaseResults, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		if err != nil {
			panic(err.Error())
		}
		exist = true
	}

	return exist
}

func createTrip(newTrip Trip, res http.ResponseWriter, req *http.Request) {

	exist := isTripExist(newTrip.TripID)

	if exist {
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte("409 - Duplicate trip ID"))
	} else {
		availableDriver := "NULL"
		tripProgress := 1
		if driver, err := retrieveAvailableDriver(); err == nil {
			availableDriver = driver.DriverID
			tripProgress = 2
		}

		query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%d')", newTrip.TripID, newTrip.PassengerID, availableDriver, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, tripProgress, currentMs(), 0)

		_, err := db.Query(query)

		if err != nil {
			panic(err.Error())
		}

		updatePassengerAvailableStatus(2, newTrip.PassengerID)
		if availableDriver != "NULL" {
			updateDriverAvailableStatus(2, availableDriver)
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("201 - Trip added: " + newTrip.TripID))
	}
}

func retrieveAvailableDriver() (*Driver, error) {
	var result []Driver

	url := "http://localhost:5000/api/v1/drivers?available_status=1"

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + authenticationToken

	// Create a new request using http
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if len(result) == 0 {
		return nil, errors.New("no available drivers are found")
	}

	return &result[0], nil
}

func updateDriverAvailableStatus(availableStatus int, driverID string) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + authenticationToken

	body := make(map[string]interface{})

	body["driver_id"] = driverID
	body["available_status"] = availableStatus

	// initialize http client
	client := &http.Client{}

	// marshal User to json
	json, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5000/api/v1/drivers/%s", driverID), bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}

func updatePassengerAvailableStatus(availableStatus int, passengerID string) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + authenticationToken

	body := make(map[string]interface{})

	body["passenger_id"] = passengerID
	body["available_status"] = availableStatus

	// initialize http client
	client := &http.Client{}

	// marshal User to json
	json, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5000/api/v1/passengers/%s", passengerID), bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}

func updateTripCompletedTime(tripid string) {
	query := fmt.Sprintf("UPDATE Trips SET CompletedTime='%d' WHERE TripID='%s'", currentMs(), tripid)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func formattedUpdateTripQueryField(trip Trip) string {
	var fields []string

	if trip.DriverID != "" {
		fields = append(fields, fmt.Sprintf("DriverID='%s'", trip.DriverID))
	}

	if trip.TripProgress != 0 {
		fields = append(fields, fmt.Sprintf("TripProgress='%d'", trip.TripProgress))
	}

	return strings.Join(fields, ", ")
}

func isTripJsonCompleted(trip Trip) bool {
	tripID := strings.TrimSpace(trip.TripID)
	passengerID := strings.TrimSpace(trip.PassengerID)
	pickupPostalCode := strings.TrimSpace(trip.PickupPostalCode)
	dropoffPostalCode := strings.TrimSpace(trip.DropoffPostalCode)

	return tripID != "" && passengerID != "" && pickupPostalCode != "" && dropoffPostalCode != ""
}

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/api/v1/trips", trips).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods("GET", "PUT", "POST")

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Trips database server -- Listening at port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler))
}
