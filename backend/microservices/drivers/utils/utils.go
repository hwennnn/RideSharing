package utils

import (
	"fmt"
	"strings"

	models "drivers/models"
)

// This method is to return boolean value whether the given driver information is completed
func IsDriverJsonCompleted(driver models.Driver) bool {
	driverID := strings.TrimSpace(driver.DriverID)
	firstName := strings.TrimSpace(driver.FirstName)
	lastName := strings.TrimSpace(driver.LastName)
	mobileNumber := strings.TrimSpace(driver.MobileNumber)
	emailAddress := strings.TrimSpace(driver.EmailAddress)
	identificationNumber := strings.TrimSpace(driver.IdentificationNumber)
	carLicenseNumber := strings.TrimSpace(driver.CarLicenseNumber)

	return driverID != "" && firstName != "" && lastName != "" && mobileNumber != "" && emailAddress != "" && identificationNumber != "" && carLicenseNumber != ""
}

// This method is to convert the field query from request query parameters,
// to the sql syntax code
func FormmatedUpdateDriverQueryField(newDriver models.Driver) string {
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

// This method is to convert the field query from request query parameters,
// to the sql syntax code
func FormattedDriverQueryField(availableStatus []string) string {
	var results string

	if len(availableStatus) > 0 && availableStatus[0] != "" {
		results += fmt.Sprintf("AvailableStatus = '%s'", availableStatus[0])
	}

	if results == "" {
		return ""
	}

	return "WHERE " + results
}
