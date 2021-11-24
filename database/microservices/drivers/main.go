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

var db *sql.DB

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(res, req)
	})
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

func main() {

	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, _ = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharing")

	fmt.Println("Database opened")

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/api/v1/drivers/", drivers).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{driverid}", driver).Methods("GET", "PUT", "POST")

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Driver database server -- Listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
