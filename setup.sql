DROP TABLE IF EXISTS RideSharing


CREATE database RideSharing;

USE RideSharing;

CREATE TABLE Passenger
(
    PassengerID BIGINT NOT NULL AUTO_INCREMENT,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    CONSTRAINT PK_Passenger PRIMARY KEY (PassengerID)
);

CREATE TABLE Driver
(
    DriverID BIGINT NOT NULL AUTO_INCREMENT,
    FirstName VARCHAR(30) NOT NULL,
    LastName VARCHAR(30) NOT NULL,
    MobileNumber VARCHAR(15) NOT NULL,
    EmailAddress VARCHAR(30) NOT NUll,
    IdentificationNumber VARCHAR(30) NOT NULL,
    CarLicenseNumber VARCHAR(30) NOT NULL,
    AvailableStatus TINYINT(1) NOT NULL,
    CONSTRAINT PK_Driver PRIMARY KEY (DriverID)
);

CREATE TABLE Trip
(
    TripID BIGINT NOT NULL AUTO_INCREMENT,
    PassengerID BIGINT NULL,
    DriverID BIGINT NULL,
    PickupPostalCode VARCHAR(10) NOT NULL,
    DropoffPostalCode VARCHAR(10) NOT NULL,
    TripProgress TINYINT(1) NOT NULL,
    CONSTRAINT PK_Trip PRIMARY KEY (TripID),
    CONSTRAINT FK_Trip_Passenger FOREIGN KEY (PassengerID) REFERENCES Passenger(PassengerID),
    CONSTRAINT FK_Trip_Driver FOREIGN KEY (DriverID) REFERENCES Driver(DriverID)
);

-- Insert Passenger Data
INSERT INTO Passenger (FirstName, LastName, MobileNumber, EmailAddress) VALUES ('Hou Man', 'Wai', '6598754815', 'hwendev@gmail.com');
INSERT INTO Passenger (FirstName, LastName, MobileNumber, EmailAddress) VALUES ('Rui Quan', 'Zachary Hong', '6512345678', 'zachary@gmail.com');
INSERT INTO Passenger (FirstName, LastName, MobileNumber, EmailAddress) VALUES ('Yong Teng', 'Tee', '6511111111', 'teeyongteng@gmail.com');

-- Insert Driver Data
INSERT INTO Driver (FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('Run Lin', 'Xiong', '6522222222', 'runlin@gmail.com', 'T12345678A', 'h124j451k32jj123f', 0);
INSERT INTO Driver (FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('Zhi Quan', 'Henry Ong', '6533333333', 'henryong@gmail.com', 'T11111111C', 'dbaa541bcc85bcb3a8', 1);
INSERT INTO Driver (FirstName, LastName, MobileNumber, EmailAddress, IdentificationNumber, CarLicenseNumber, AvailableStatus) VALUES ('Ming Han', 'Vincent Tee', '6544444444', 'vincentminghan@gmail.com', 'T22222222B', 'agfahudsi142kj42', 2);

-- Insert Trip Data
INSERT INTO Trip (PickupPostalCode, DropoffPostalCode, TripProgress) VALUES ('642678', '730022', 0);
UPDATE Trip SET PassengerID=1, DriverID=2, TripProgress = 2 WHERE TripID = 1;
INSERT INTO Trip (PickupPostalCode, DropoffPostalCode, TripProgress) VALUES ('111111', '222222', 1);
INSERT INTO Trip (PickupPostalCode, DropoffPostalCode, TripProgress) VALUES ('333333', '444444', 2);



