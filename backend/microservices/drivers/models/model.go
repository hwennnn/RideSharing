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
	// Notes for Available Status
	// 0 -> used by golang to indicate whether the integer variable has been initialised or not
	// 2 => Online and assigned a trip by system (but the trip has not been initiated yet)
	// 3 => Online and during the trip
}
