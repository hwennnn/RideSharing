package main

import (
	"database/sql"
	"encoding/json"
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

func fetchDriver(url string) Driver {
	var result Driver

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(body, &result)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return result
}

func fetchPassenger(url string) Passenger {
	var result Passenger

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(body, &result)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return result
}

func trips(res http.ResponseWriter, req *http.Request) {
	var results []Trip

	params := req.URL.Query()

	query := "SELECT * FROM Trips"
	formmatedFieldQuery := formmatedTripQueryField(params["driver_id"], params["passenger_id"], params["trip_progress"])
	if formmatedFieldQuery != "" {
		query = fmt.Sprintf("SELECT * FROM Trips WHERE %s", formmatedFieldQuery)
	}
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
			trip.Driver = fetchDriver(fmt.Sprintf("http://localhost:5000/api/v1/drivers/%s", trip.DriverID))
		}

		if trip.PassengerID != "" {
			trip.Passenger = fetchPassenger(fmt.Sprintf("http://localhost:5000/api/v1/passengers/%s", trip.PassengerID))
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
	fmt.Println(driverID, passengerID)

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

	return results
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
				trip.Driver = fetchDriver(fmt.Sprintf("http://localhost:5000/api/v1/drivers/%s", trip.DriverID))
			}

			if trip.PassengerID != "" {
				trip.Passenger = fetchPassenger(fmt.Sprintf("http://localhost:5000/api/v1/passengers/%s", trip.PassengerID))
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
					query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', NULL, '%s', '%s', '%d', '%d', '%d')", newTrip.TripID, newTrip.PassengerID, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, 1, currentMs(), 0)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}

					updatePassengerAvailableStatus(2, newTrip.PassengerID)

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

				if tripFromDatabase.TripProgress != 2 && newTrip.TripProgress == 2 {
					// update driver and passenger available status back to 1 once the trip is completed
					updateDriverAvailableStatus(2, newTrip.DriverID)
				} else if tripFromDatabase.TripProgress != 3 && newTrip.TripProgress == 3 {
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

func updateDriverAvailableStatus(availableStatus int, driverID string) {
	query := fmt.Sprintf("UPDATE Drivers SET AvailableStatus='%d' WHERE DriverID='%s'", availableStatus, driverID)
	fmt.Println(query)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func updatePassengerAvailableStatus(availableStatus int, passengerID string) {
	query := fmt.Sprintf("UPDATE Passengers SET AvailableStatus='%d' WHERE PassengerID='%s'", availableStatus, passengerID)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
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
