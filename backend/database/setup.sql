/* CREATE database RideSharing; */

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

-- Insert Passenger Data
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress, AvailableStatus) VALUES ('1', 'Hou Man', 'Wai', '6598754815', 'hwendev@gmail.com', 1);
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress, AvailableStatus) VALUES ('2' ,'Rui Quan', 'Zachary Hong', '6512345678', 'zachary@gmail.com', 1);
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress, AvailableStatus) VALUES ('3', 'Yong Teng', 'Tee', '6511111111', 'teeyongteng@gmail.com', 1);

-- Insert Driver Data
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('1', 'Run Lin', 'Xiong', '6522222222', 'runlin@gmail.com', 'T12345678A', 'h124j451k32jj123f', 0);
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('2', 'Zhi Quan', 'Henry Ong', '6533333333', 'henryong@gmail.com', 'T11111111C', 'dbaa541bcc85bcb3a8', 1);
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('3', 'Ming Han', 'Vincent Tee', '6544444444', 'vincentminghan@gmail.com', 'T22222222B', 'agfahudsi142kj42', 2);

-- Insert Trip Data
INSERT INTO Trips (TripID, PassengerID, DriverID, PickupPostalCode, DropoffPostalCode, TripProgress, CreatedTime, CompletedTime) VALUES ('1', '1', '2', '642678', '730022', 1, 1637424569024, 0);
INSERT INTO Trips (TripID, PassengerID, DriverID, PickupPostalCode, DropoffPostalCode, TripProgress, CreatedTime, CompletedTime) VALUES ('2', '2', '3', '111111', '222222', 1, 1637524569024, 0);
INSERT INTO Trips (TripID, PassengerID, DriverID, PickupPostalCode, DropoffPostalCode, TripProgress, CreatedTime, CompletedTime) VALUES ('3', '1', '3', '333333', '444444', 1, 1637424565554, 0);
