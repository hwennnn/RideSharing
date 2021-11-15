DROP TABLE IF EXISTS RideSharing


CREATE database RideSharing;

USE RideSharing;

CREATE TABLE Passenger
(
    PassengerID INTEGER NOT NULL AUTO_INCREMENT,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    CONSTRAINT PK_Passenger PRIMARY KEY (PassengerID)
);

CREATE TABLE Driver
(
    DriverID INTEGER NOT NULL AUTO_INCREMENT,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    IdentificationNumber VARCHAR(30) NOT NULL,
    CarLicenseNumber VARCHAR(30) NOT NULL,
    CONSTRAINT PK_Driver PRIMARY KEY (DriverID)
);

CREATE TABLE Trip
(
    TripID INTEGER NOT NULL AUTO_INCREMENT,
    PassengerID INTEGER NULL,
    DriverID INTEGER NULL,
    PickupPostalCode VARCHAR(10) NOT NULL,
    DropoffPostalCOde VARCHAR(10) NOT NULL,
    Progress INTEGER NOT NULL,
    CONSTRAINT PK_Trip PRIMARY KEY (TripID),
    CONSTRAINT FK_Trip_Passenger FOREIGN KEY (PassengerID) REFERENCES Passenger(PassengerID),
    CONSTRAINT FK_Trip_Driver FOREIGN KEY (DriverID) REFERENCES Driver(DriverID)
);