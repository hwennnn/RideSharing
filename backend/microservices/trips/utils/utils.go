package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	models "trips/models"

	_ "github.com/go-sql-driver/mysql"
)

const authenticationToken = "2a6b36bf-61b9-4d0e-904c-7843e7b97308"
const serverEndpointBaseURL = "http://acl:4000/api/v1"

var driverEndpointBaseURL = fmt.Sprintf("%s/drivers", serverEndpointBaseURL)
var passengerEndpointBaseURL = fmt.Sprintf("%s/passengers", serverEndpointBaseURL)

// retrieve current milliseconds since epoch
func CurrentMs() int64 {
	return time.Now().Round(time.Millisecond).UnixNano() / 1e6
}

// This method is to convert the field query from request query parameters,
// to the sql syntax code
func FormmatedTripQueryField(driverID []string, passengerID []string, tripProgress []string) string {
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

// This method is to convert the field query from request query parameters,
// to the sql syntax code
func FormattedUpdateTripQueryField(trip models.Trip) string {
	var fields []string

	if trip.DriverID != "" {
		fields = append(fields, fmt.Sprintf("DriverID='%s'", trip.DriverID))
	}

	if trip.TripProgress != 0 {
		fields = append(fields, fmt.Sprintf("TripProgress='%d'", trip.TripProgress))
	}

	return strings.Join(fields, ", ")
}

// This method is to return boolean value whether the given trip information is completed
func IsTripJsonCompleted(trip models.Trip) bool {
	tripID := strings.TrimSpace(trip.TripID)
	passengerID := strings.TrimSpace(trip.PassengerID)
	pickupPostalCode := strings.TrimSpace(trip.PickupPostalCode)
	dropoffPostalCode := strings.TrimSpace(trip.DropoffPostalCode)

	return tripID != "" && passengerID != "" && pickupPostalCode != "" && dropoffPostalCode != ""
}

// This method will send a request to driver microservice
// in order to retrieve driver information
func FetchDriver(driverID string) models.Driver {
	var result models.Driver

	url := fmt.Sprintf("%s/%s", driverEndpointBaseURL, driverID)

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

// This method will send a request to passenger microservice
// in order to retrieve passenger information
func FetchPassenger(passengerID string) models.Passenger {
	var result models.Passenger

	url := fmt.Sprintf("%s/%s", passengerEndpointBaseURL, passengerID)

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

// This method will send a request to driver microservice
// in order to retrieve available drivers information
func RetrieveAvailableDriver() (*models.Driver, error) {
	var result []models.Driver

	url := fmt.Sprintf("%s?available_status=1", driverEndpointBaseURL)

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

// This method will send a request to driver microservice
// in order to update driver available status
func UpdateDriverAvailableStatus(availableStatus int, driverID string) {
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
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", driverEndpointBaseURL, driverID), bytes.NewBuffer(json))
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

// This method will send a request to passenger microservice
// in order to update passenger available status
func UpdatePassengerAvailableStatus(availableStatus int, passengerID string) {
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
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", passengerEndpointBaseURL, passengerID), bytes.NewBuffer(json))
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
