package main

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	utils "backend/microservices/trips/utils"
	models "backend/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// global database handler object
var db *sql.DB

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
}

func getTrips(res http.ResponseWriter, req *http.Request) {
	var results []models.Trip

	params := req.URL.Query()

	formmatedFieldQuery := utils.FormmatedTripQueryField(params["driver_id"], params["passenger_id"], params["trip_progress"])
	query := fmt.Sprintf("SELECT * FROM Trips %s ORDER BY CompletedTime DESC", formmatedFieldQuery)
	databaseResults, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		var trip models.Trip
		var sqlDriverID sql.NullString
		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &sqlDriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime)
		if sqlDriverID.Valid {
			trip.DriverID = sqlDriverID.String
		}

		if trip.DriverID != "" {
			trip.Driver = utils.FetchDriver(trip.DriverID)
		}

		if trip.PassengerID != "" {
			trip.Passenger = utils.FetchPassenger(trip.PassengerID)
		}

		if err != nil {
			panic(err.Error())
		}
		results = append(results, trip)

	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func getTrip(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	tripid := params["tripid"]

	isDriverExist, trip := getTripHelper(tripid)

	if isDriverExist {
		json.NewEncoder(res).Encode(trip)
	} else {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("404 - No Trip found"))
	}
}

// check if trip exists by tripid
func getTripHelper(tripid string) (bool, models.Trip) {
	query := fmt.Sprintf("SELECT * FROM Trips WHERE TripID='%s'", tripid)
	databaseResults, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	var isExist bool
	var trip models.Trip
	var sqlDriverID sql.NullString
	for databaseResults.Next() {
		err = databaseResults.Scan(&trip.TripID, &trip.PassengerID, &sqlDriverID, &trip.PickupPostalCode, &trip.DropoffPostalCode, &trip.TripProgress, &trip.CreatedTime, &trip.CompletedTime)

		if sqlDriverID.Valid {
			trip.DriverID = sqlDriverID.String
		}

		if trip.DriverID != "" {
			trip.Driver = utils.FetchDriver(trip.DriverID)
		}

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

		isTripExist, _ := getTripHelper(newTrip.TripID)

		if isTripExist {
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte("409 - Duplicate trip ID"))
		} else {
			availableDriver := "NULL"
			tripProgress := 1

			if driver, err := utils.RetrieveAvailableDriver(); err == nil {
				availableDriver = driver.DriverID
				tripProgress = 2
			}

			query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%d')", newTrip.TripID, newTrip.PassengerID, availableDriver, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, tripProgress, utils.CurrentMs(), 0)

			_, err := db.Query(query)

			if err != nil {
				panic(err.Error())
			}

			utils.UpdatePassengerAvailableStatus(2, newTrip.PassengerID)
			if availableDriver != "NULL" {
				utils.UpdateDriverAvailableStatus(2, availableDriver)
			}

			res.WriteHeader(http.StatusCreated)
			res.Write([]byte("201 - Trip added: " + newTrip.TripID))
		}

	} else {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write([]byte("422 - Please supply trip information in JSON format"))
	}
}

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

			availableDriver := "NULL"
			tripProgress := 1

			if driver, err := utils.RetrieveAvailableDriver(); err == nil {
				availableDriver = driver.DriverID
				tripProgress = 2
			}

			query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%d')", newTrip.TripID, newTrip.PassengerID, availableDriver, newTrip.PickupPostalCode, newTrip.DropoffPostalCode, tripProgress, utils.CurrentMs(), 0)

			_, err := db.Query(query)

			if err != nil {
				panic(err.Error())
			}

			utils.UpdatePassengerAvailableStatus(2, newTrip.PassengerID)
			if availableDriver != "NULL" {
				utils.UpdateDriverAvailableStatus(2, availableDriver)
			}

			res.WriteHeader(http.StatusCreated)
			res.Write([]byte("201 - Trip added: " + newTrip.TripID))
		} else {
			formattedUpdateFieldQuery := utils.FormattedUpdateTripQueryField(newTrip)

			if formattedUpdateFieldQuery == "" { // means there is no valid field can be updated
				res.WriteHeader(http.StatusUnprocessableEntity)
				res.Write([]byte("422 - Please supply trip information in JSON format"))
				return
			}

			if tripFromDatabase.TripProgress != 3 && newTrip.TripProgress == 3 {
				// update driver available status to 3 as driver has initiated the trip and currently during the trip
				utils.UpdateDriverAvailableStatus(3, newTrip.DriverID)
			} else if tripFromDatabase.TripProgress != 4 && newTrip.TripProgress == 4 {
				// update driver and passenger available status back to 1 once the trip is completed
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

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Trips database server -- Listening at port 8082")
	log.Fatal(http.ListenAndServe(":8082", handler))
}
