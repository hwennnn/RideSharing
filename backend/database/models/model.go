package models

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
	// 1 -> Created by passenger, but no driver is found to be assgined
	// 2 -> A driver was already assigned for the trip, but the driver has not inititated the trip yet
	// 3 -> The trip is ongoing (Driver has initiated the trip)
	// 4 -> The trip has ended (Driver has ended the trip)
}
