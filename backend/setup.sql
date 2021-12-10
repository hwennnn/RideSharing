CREATE database RideSharing;

USE RideSharing;

CREATE TABLE Passengers
(
    PassengerID VARCHAR(36) NOT NULL,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    AvailableStatus TINYINT(1) NOT NULL,
    CONSTRAINT PK_Passenger PRIMARY KEY (PassengerID)
);

CREATE TABLE Drivers
(
    DriverID VARCHAR(36) NOT NULL,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    IdentificationNumber VARCHAR(30) NOT NULL,
    CarLicenseNumber VARCHAR(30) NOT NULL,
    AvailableStatus TINYINT(1) NOT NULL,
    CONSTRAINT PK_Driver PRIMARY KEY (DriverID)
);

CREATE TABLE Trips
(
    TripID VARCHAR(36) NOT NULL,
    PassengerID VARCHAR(36) NOT NULL,
    DriverID VARCHAR(36) NULL,
    PickupPostalCode VARCHAR(10) NOT NULL,
    DropoffPostalCode VARCHAR(10) NOT NULL,
    TripProgress TINYINT(1) NOT NULL,
    CreatedTime BIGINT NOT NULL,
    CompletedTime BIGINT NULL,
    CONSTRAINT PK_Trip PRIMARY KEY (TripID)
);