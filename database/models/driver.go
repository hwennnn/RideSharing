package driver

type Driver struct {
	DriverID             int64
	FirstName            string
	LastName             string
	MobileNumber         string
	EmailAddress         string
	IdentificationNumber string
	CarLicenseNumber     string
	AvailableStatus      int
	// 0 => Offline
	// 1 => Online and available
	// 2 => Online and during the trip
}
