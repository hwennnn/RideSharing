package main

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	models "trips/models"
	utils "trips/utils"

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

// This method is used to retrieve trips from MySQL,
// and return the result in array of trip json object
func getTrips(res http.ResponseWriter, req *http.Request) {
	var results []models.Trip

	params := req.URL.Query()

	// Customise the field query from request query parameters
	formmatedFieldQuery := utils.FormmatedTripQueryField(params["driver_id"], params["passenger_id"], params["trip_progress"])

	// sort the result by completed time in reverse order
	query := fmt.Sprintf("SELECT * FROM Trips %s ORDER BY CompletedTime DESC", formmatedFieldQuery)
	databaseResults, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		var trip models.Trip

		// use sql null string interface to check whether a value is provided by the user
		var sqlDriverID sql.NullString

		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &sqlDriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime)

		// if the driverID is valid, update in the trip object
		if sqlDriverID.Valid {
			trip.DriverID = sqlDriverID.String
		}

		// retrieve the driver information by sending http request to driver microservice
		if trip.DriverID != "" {
			trip.Driver = utils.FetchDriver(trip.DriverID)
		}

		// retrieve the driver information by sending http request to passenger microservice
		if trip.PassengerID != "" {
			trip.Passenger = utils.FetchPassenger(trip.PassengerID)
		}

		if err != nil {
			panic(err.Error())
		}

		results = append(results, trip)
	}

	// returns all the trips in JSON
	json.NewEncoder(res).Encode(results)
}

// This method is used to retrieve a trip from MySQL by specific tripID,
// and return the result in json otherwise return 404 code
func getTrip(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	tripid := params["tripid"]

	isTripExist, trip := getTripHelper(tripid)

	if isTripExist {
		json.NewEncoder(res).Encode(trip)
	} else {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("404 - No Trip found"))
	}
}

// This helper method helps to query the trip from the database,
// and return (boolean, trip) tuple object
func getTripHelper(tripid string) (bool, models.Trip) {
	query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID='%s'", tripid)
	databaseResults, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	var isExist bool
	var trip models.Trip

	// use sql null string interface to check whether a value is provided by the user
	var sqlDriverID sql.NullString

	for databaseResults.Next() {
		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &sqlDriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime)

		// if the driverID is valid, update in the trip object
		if sqlDriverID.Valid {
			trip.DriverID = sqlDriverID.String
		}

		// retrieve the driver information by sending http request to driver microservice
		if trip.DriverID != "" {
			trip.Driver = utils.FetchDriver(trip.DriverID)
		}

		// retrieve the driver information by sending http request to passenger microservice
		if trip.PassengerID != "" {
			trip.Passenger = utils.FetchPassenger(trip.PassengerID)
		}

		if err != nil {
			panic(err.Error())
		}
		isExist = true
	}

	return isExist, trip
}

// This method is used to create a trip in MySQL by specific tripID,
// Case 1: If the compulsory trip information is not provided, it will return message which says the information is not correctly supplied
// Case 2: It will fail and return conflict status code if a trip with same tripID is already found in the database
// Case 3: Otherwise, it will return success message with status created code
func postTrip(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	tripid := params["tripid"]

	// read the string sent to the service
	var newTrip models.Trip
	reqBody, err := ioutil.ReadAll(req.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &newTrip)

		if !utils.IsTripJsonCompleted(newTrip) {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Please supply trip information in JSON format"))
			return
		}

		if tripid != newTrip.TripID {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - The data in body and parameters do not match"))
			return
		}

		// check if trip exists; add only if trip does not exist
		isTripExist, _ := getTripHelper(newTrip.TripID)

		if isTripExist {
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte("409 - Duplicate trip ID"))
		} else {
			createTripHelper(newTrip)

			res.WriteHeader(http.StatusCreated)
			res.Write([]byte("201 - Trip added: " + newTrip.TripID))
		}

	} else {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte("422 - Please supply trip information in JSON format"))
	}
}

// This method is used to insert a new trip into the database
func createTripHelper(newTrip models.Trip) {
	// initially set the available driver to null first (the driverID is nullable in the database)
	availableDriver := "NULL"
	// initially set the progress as trip is not assigned a driver yet
	tripProgress := 1

	// if an available driver is found,
	// update the availableDriver id and trip progress respectively
	if driver, err := utils.RetrieveAvailableDriver(); err == nil {
		availableDriver = driver.DriverID
		tripProgress = 2
	}

	query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%d')", newTrip.TripID, newTrip.PassengerID, availableDriver, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, tripProgress, utils.CurrentMs(), 0)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	// Since the trip has been created (meaning the trip is in the progress)
	// update the passenger available status to 2
	// update the driver avalable status to 2 if there is available driver
	utils.UpdatePassengerAvailableStatus(2, newTrip.PassengerID)
	if availableDriver != "NULL" {
		utils.UpdateDriverAvailableStatus(2, availableDriver)
	}
}

// This method is used for either creating or updating the trip depends whether the tripID exists
// Case 1: If tirpID exists, update the trip using the information retrieved from request body
// Case 2: If tripID does not exist, create the trip using the information retrieved from request body
func putTrip(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	tripid := params["tripid"]

	var newTrip models.Trip
	reqBody, err := ioutil.ReadAll(req.Body)

	if err == nil {
		json.Unmarshal(reqBody, &newTrip)

		if tripid != newTrip.TripID {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - The data in body and parameters do not match"))
			return
		}

		// check if trip exists; add only if trip does not exist
		isTripExist, tripFromDatabase := getTripHelper(tripid)

		if !isTripExist {
			if !utils.IsTripJsonCompleted(newTrip) {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply trip information in JSON format"))
				return
			}

			createTripHelper(newTrip)

			res.WriteHeader(http.StatusCreated)
			res.Write([]byte("201 - Trip added: " + newTrip.TripID))
		} else {
			formattedUpdateFieldQuery := utils.FormattedUpdateTripQueryField(newTrip)

			// means there is no valid field can be updated
			if formattedUpdateFieldQuery == "" {
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply trip information in JSON format"))
				return
			}

			if tripFromDatabase.TripProgress != 3 && newTrip.TripProgress == 3 {
				// As driver has initiated the trip and currently during the trip
				// Update driver available status to 3
				utils.UpdateDriverAvailableStatus(3, newTrip.DriverID)
			} else if tripFromDatabase.TripProgress != 4 && newTrip.TripProgress == 4 {
				// Once the trip is completed
				// Update driver and passenger available status back to 1
				// Update trip completed time as current milliseconds since epoch
				utils.UpdateDriverAvailableStatus(1, tripFromDatabase.DriverID)
				utils.UpdatePassengerAvailableStatus(1, tripFromDatabase.PassengerID)
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

// This method is used to update trip completed time as current milliseconds since epoch
func updateTripCompletedTime(tripid string) {
	query := fmt.Sprintf("UPDATE Trips SET CompletedTime='%d' WHERE TripID='%s'", utils.CurrentMs(), tripid)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/api/v1/trips", getTrips).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripid}", getTrip).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripid}", postTrip).Methods("POST")
	router.HandleFunc("/api/v1/trips/{tripid}", putTrip).Methods("PUT")

	// enable cross-origin resource sharing (cors) for all requests
	handler := cors.AllowAll().Handler(router)

	fmt.Println("Trips database server -- Listening at port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler))
}
