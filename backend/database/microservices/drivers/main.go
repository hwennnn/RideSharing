package main

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	utils "backend/microservices/drivers/utils"
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

func getDrivers(res http.ResponseWriter, req *http.Request) {
	var results []models.Driver

	params := req.URL.Query()

	formmatedFieldQuery := utils.FormattedDriverQueryField(params["available_status"])
	query := fmt.Sprintf("SELECT * FROM Drivers %s", formmatedFieldQuery)

	databaseResults, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for databaseResults.Next() {
		// map this type to the record in the table
		var driver models.Driver
		err = databaseResults.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.IdentificationNumber, &driver.CarLicenseNumber, &driver.AvailableStatus)
		if err != nil {
			panic(err.Error())
		}
		results = append(results, driver)
	}
	// returns all the courses in JSON
	json.NewEncoder(res).Encode(results)
}

func getDriver(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	driverid := params["driverid"]

	query := fmt.Sprintf("SELECT * FROM Drivers WHERE DriverID='%s'", driverid)
	databaseResults, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	var isExist bool
	var driver models.Driver
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

func postDriver(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	driverid := params["driverid"]

	// read the string sent to the service
	var newDriver models.Driver
	reqBody, err := ioutil.ReadAll(req.Body)

	if err == nil {
		// convert JSON to object
		json.Unmarshal(reqBody, &newDriver)

		if !utils.IsDriverJsonCompleted(newDriver) {
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

func putDriver(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	driverid := params["driverid"]

	var newDriver models.Driver
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
			if !utils.IsDriverJsonCompleted(newDriver) {
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
			formattedUpdateFieldQuery := utils.FormmatedUpdateDriverQueryField(newDriver)

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

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/api/v1/drivers", getDrivers).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{driverid}", getDriver).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{driverid}", postDriver).Methods("POST")
	router.HandleFunc("/api/v1/drivers/{driverid}", putDriver).Methods("PUT")

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Driver database server -- Listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
