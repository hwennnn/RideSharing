package utils

import (
	"fmt"
	"strings"

	models "passengers/models"
)

// This method is to convert the field query from request query parameters,
// to the sql syntax code
func FormmatedUpdatePassengerQueryField(newPassenger models.Passenger) string {
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

// This method is to return boolean value whether the given passenger information is completed
func IsPassengerJsonCompleted(passenger models.Passenger) bool {
	passengerID := strings.TrimSpace(passenger.PassengerID)
	firstName := strings.TrimSpace(passenger.FirstName)
	lastName := strings.TrimSpace(passenger.LastName)
	mobileNumber := strings.TrimSpace(passenger.MobileNumber)
	emailAddress := strings.TrimSpace(passenger.EmailAddress)

	return passengerID != "" && firstName != "" && lastName != "" && mobileNumber != "" && emailAddress != ""
}
