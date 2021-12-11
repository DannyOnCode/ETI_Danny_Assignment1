DROP database DRidePassenger;
DROP database DRideDriver;
DROP database DRideTrip;

CREATE database DRidePassenger;
CREATE database DRideDriver;
CREATE database DRideTrip;

USE DRidePassenger;
CREATE TABLE Passenger (PassengerID VARCHAR(100) NOT NULL PRIMARY KEY, FirstName VARCHAR(30), LastName VARCHAR(30), MobileNo CHAR(8), Email VARCHAR(100)); 

USE DRideDriver;
CREATE TABLE Driver (DriverID VARCHAR(100) NOT NULL PRIMARY KEY, FirstName VARCHAR(30), LastName VARCHAR(30), MobileNo CHAR(8), Email VARCHAR(100), LicenseNumber VARCHAR(100), Status VARCHAR(100));

USE DRideTrip;
CREATE TABLE Trip (TripID int NOT NULL AUTO_INCREMENT, PassengerID VARCHAR(30), DriverID VARCHAR(30), PickUpPostalCode VARCHAR(10), DropOffPostalCode VARCHAR(10), StartDateTime datetime, EndDateTime datetime, PRIMARY KEY (TripID));

USE DRideDriver;
INSERT INTO Driver
    VALUES ('D1', 'Danny', 'Chan', 91111111, 'nihility@gmail.com', 'SBA1234D', 'Available');
INSERT INTO Driver
    VALUES ('D2', 'Oh Hak', 'Eews', 92222222, 'tokyodriftnunu@gmail.com', 'SJT9876K', 'Available');
USE DRidePassenger;
INSERT INTO Passenger
    VALUES ('P1', 'Kenneth', 'Back', 93333333, 'backer@gmail.com');
INSERT INTO Passenger
    VALUES ('P2', 'Pritheev', 'Oofer', 94444444, 'PrettyOof@gmail.com');