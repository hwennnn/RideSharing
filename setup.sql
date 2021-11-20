CREATE database RideSharing;

USE RideSharing;

CREATE TABLE Passengers
(
    PassengerID VARCHAR(30) NOT NULL,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    CONSTRAINT PK_Passenger PRIMARY KEY (PassengerID)
);

CREATE TABLE Drivers
(
    DriverID VARCHAR(30) NOT NULL,
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
    TripID VARCHAR(30) NOT NULL,
    PassengerID VARCHAR(30) NOT NULL,
    DriverID VARCHAR(30) NOT NULL,
    PickupPostalCode VARCHAR(10) NOT NULL,
    DropoffPostalCode VARCHAR(10) NOT NULL,
    TripProgress TINYINT(1) NOT NULL,
    CreatedTime BIGINT NOT NULL,
    CompletedTime BIGINT NULL,
    CONSTRAINT PK_Trip PRIMARY KEY (TripID),
    CONSTRAINT FK_Trip_Passenger FOREIGN KEY (PassengerID) REFERENCES Passengers(PassengerID),
    CONSTRAINT FK_Trip_Driver FOREIGN KEY (DriverID) REFERENCES Drivers(DriverID)
);

-- Insert Passenger Data
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress) VALUES ('1', 'Hou Man', 'Wai', '6598754815', 'hwendev@gmail.com');
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress) VALUES ('2' ,'Rui Quan', 'Zachary Hong', '6512345678', 'zachary@gmail.com');
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNumber, EmailAddress) VALUES ('3', 'Yong Teng', 'Tee', '6511111111', 'teeyongteng@gmail.com');

-- Insert Driver Data
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('1', 'Run Lin', 'Xiong', '6522222222', 'runlin@gmail.com', 'T12345678A', 'h124j451k32jj123f', 0);
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('2', 'Zhi Quan', 'Henry Ong', '6533333333', 'henryong@gmail.com', 'T11111111C', 'dbaa541bcc85bcb3a8', 1);
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('3', 'Ming Han', 'Vincent Tee', '6544444444', 'vincentminghan@gmail.com', 'T22222222B', 'agfahudsi142kj42', 2);

-- Insert Trip Data
INSERT INTO Trips (TripID, PassengerID, DriverID, PickupPostalCode, DropoffPostalCode, TripProgress, CreatedTime) VALUES ('1', '1', '2', '642678', '730022', 1, 1637424569024);
UPDATE Trips SET TripProgress = 2 WHERE TripID = '1';
INSERT INTO Trips (TripID, PassengerID, DriverID, PickupPostalCode, DropoffPostalCode, TripProgress, CreatedTime) VALUES ('2', '2', '3', '111111', '222222', 1, 1637524569024);
UPDATE Trips SET TripProgress = 3 WHERE TripID = '2';
INSERT INTO Trips (TripID, PassengerID, DriverID, PickupPostalCode, DropoffPostalCode, TripProgress, CreatedTime) VALUES ('3', '1', '3', '333333', '444444', 2, 1637424565554);
UPDATE Trips SET TripProgress = 2 WHERE TripID = '3';


SELECT * FROM Trips t
INNER JOIN Drivers d ON t.DriverID = d.DriverID 
INNER JOIN Passengers p ON t.PassengerID = p.PassengerID
WHERE t.DriverID = '2'