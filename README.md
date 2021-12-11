## Table of Contents
1. [Introduction](#introduction)
2. [Design Consideration of Microservices](#consideration)
3. [Microservice APIs of Application](#microservices_explain)
    1. [Passenger](#passenger)
    2. [Driver](#driver)
    3. [Trip](#trip)
4. [Instructions to setting up and running of microservices](#instructions)
    1. [Database](#database)
    2. [Microservices](#microservice)

## Introduction  <a name="introduction"></a>
I am the developer Danny Chan Yu Tian and This README is dedicated to talk about the follow assignment, ETI Assignment 1 on a Trip Application called "DRide".<br> 
In the follow sections I will first be discussing the design consideration that I have made when designing the application and its microservices.<br>
I will then show a diagram of the architecture of application and be explaining each Microservices.<br>
Lastly, I will end the README with instructions on how to set up and run the microservice and database<br>

## Design Consideration of Microservices <a name="consideration"></a>
DRide will consist of a web service (a frontend), 3 microservices as well as 3 databases, each connected to one microservice.<br>

Below is an image of the design of the entire application <br>

![Architecture of Application](/images/Architecture.png) <br>

Firstly, I will quickly talk about the design flow. The user will interact with the front end and the front end will send requests to the APIs depending on the needs. <br>The Microservice will then process this requests and access the databases whenever needed. The data will then be send back to the front end to be displayed.<br>

I will now go through some of the consideration that I have made when implementing each of these microservice and front end.<br>
When I was considering which front end service to use, I initially considered using REACT as my framework however after doing more research I decided to use html template to serve my html pages as all of my pages are static pages. This method is also the easiest and simpliest to get the result that I need.

For the backend, I have implemented 3 different microswervices of passenger, driver and trip as seen in the diagram. Initially I connected all three microservice to one database however after careful consideration I have decided to have one database for each microservice to allow for a more loosely coupled application. By doing so, if one database is down, the other functions will not be affected. Besides that, these microservices will also be hosted on different ports, Passenger on :80 Driver on :100, and Trip on :120.

## Microservice APIs of Application <a name="microservices_explain"></a>

Under this section, I will be discussing each microservice and the resources that they provide along with the routes to access the resources.
### Passenger API <a name="passenger"></a>
| Request   | URL |
|--------   |------|
| GET       | /api/v1/passenger/{passengerID}?mobileNo={passengerNumber} |
| POST      | /api/v1/passenger/{passengerID} |
| PUT       | /api/v1/passenger/{passengerID}?mobileNo={passengerNumber} |
| DELETE    | /api/v1/passenger/{passengerID} |

I will now be describing the passenger URL requests with they purpose.<br>
Firstly, the GET request is a login function where it takes in the user's ID and mobile number to run the check against the database to ensure the login data is correct. Once the data is correct, the GET request will send back the full detail of the inputted passenger id, details such as, passenger id, first name, last name, mobile number, and email. <br><br>

The POST request is used to register new passengers. When requesting this api with the post method, the passenger will also pass in the full struct with details of the passenger such as passenger id, first name, last name, mobile number, and email using json. The post method will then process this information and create a new passenger inside the database.<br><br>

The PUT request is used to update the passenger's information. When requesting this api with the PUT method, the passenger will also pass in a full struct with details of the passenger such as first name, last name, mobile number, and email using json. If any of these attributes are not filled in, the api will only update the fields that the user has filled in.<br><br>

The DELETE request is used to delete a passenger record, however as the requirements states that accounts cannot be deleted, the http request will send back a message stating 403 - unable to delete account due to auditting reasons.<br><br>

### Driver API <a name="driver"></a>
| Request   | URL |
|--------   |------|
| GET       | /api/v1/driver/{driverID}?mobileNo={driverNumber} |
| POST      | /api/v1/driver/{driverID} |
| PUT       | /api/v1/driver/{driverID}?mobileNo={driverNumber} |
| DELETE    | /api/v1/driver/{driverID} |

I will now be describing the driver URL requests with they purpose.<br>
This api call is very similar to the passenger API.<br>
Firstly, the GET request is a login function where it takes in the user's ID and mobile number to run the check against the database to ensure the login data is correct. Once the data is correct, the GET request will send back the full detail of the inputted passenger id, details such as, driver id, first name, last name, mobile number, email, and License Number. <br><br>

The POST request is used to register new drivers. When requesting this api with the post method, the driver will also pass in the full struct with details of the driver such as driver id, first name, last name, mobile number, email, and License Number using json. The post method will then process this information and create a new driver inside the database.<br><br>

The PUT request is used to update the driver's information. When requesting this api with the PUT method, the driver will also pass in a full struct with details of the driver such as first name, last name, mobile number,email, and License Number using json. If any of these attributes are not filled in, the api will only update the fields that the user has filled in.<br><br>

The DELETE request is used to delete a driver record, however as the requirements states that accounts cannot be deleted, the http request will send back a message stating 403 - unable to delete account due to auditting reasons.<br><br>

### Trip API <a name="trip"></a>
| Request   | URL  |
|-----------|-------|
| GET       | /api/v1/trip/{ID}?userType={usertype}      |
| POST      | /api/v1/trip/{ID}      |
| PUT       | /api/v1/trip/{ID}      |

I will now be describing the trip URL requests with they purpose.<br>
Firstly, the get request takes in 3 different type of userType, Driver, Passenger or nil. When the user type is driver, the GET request will return any trips that the driver may have been assigned to, this allows the driver to start a trip and end trips. If the user type is a passenger, it will retrieve all past trips that the passenger has taken and return an array of all the past trips. If the user type is not specified the API will return the trip record of the index TripID.<br><br>

The POST request is used when passenger wants to request for a trip. When requesting this api with the POST method, the passenger will also pass into the request the Pick up and Drop off postal codes using json. Once the request has been processed, the api will return the driver that has been assigned to the trip.<br><br>

The PUT request is used by the driver that wants to either start a trip or end a trip. When requested, the API will check whether the Trip that was indicated is a ongoing trip or a new trip. This is done by checking the StartDateTime and EndDateTime. If the StartDateTime is null then the API will start the trip and if the EndDateTime is null but the StartDateTime is not, then the API will end the trip. These information will then be updated to the database.<br><br>


## Instructions to setting up and running of microservices <a name="instructions"></a>
Under this section, I will be going through how to start up the microservices that have been developed as well as how to set up the database to create the databases and tables.
### Setting up of Database <a name="database"></a>
Under the "ETI_DANNY_ASSIGNMENT1" Folder there is a file named "sql_eti_setup.sql". Open the file and run the script inside the file. By doing so, you should have created the needed databases and tables required.<br>
Three database will be created called, "DRidePassenger", "DRideDriver", "DRideTrip".<br>
Under each database the tables, "Passenger", "Driver", "Trip" should have been created respectively.<br>
Lastly, some dummy data will have been created under the Passenger and Driver.<br>

### Commands to run Microservices <a name="microservice"></a>
To run the microservices, firstly open four different command prompts in your terminal.<br><br>
In the first command prompt access the web folder and start the microservice through the follow commands
```console
cd web && go run main.go
```
In the second command prompt access the passenger folder and start the microservice through the follow commands
```console
cd passenger && go run main.go
```
In the third command prompt access the driver folder and start the microservice through the follow commands
```console
cd driver && go run main.go
```
In the second command prompt access the trip folder and start the microservice through the follow commands
```console
cd trip && go run main.go
```

