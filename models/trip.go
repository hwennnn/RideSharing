package trip

type Trip struct {
	TripID            int64
	PassengerID       int64
	DriverID          int64
	PickupPostalCode  string
	DropoffPostalCode string
	TripProgress      int
	// 0 -> Created by passenger, but no driver is assigned yet
	// 1 -> The trip is ongoing
	// 2 -> The trip has ended
}
