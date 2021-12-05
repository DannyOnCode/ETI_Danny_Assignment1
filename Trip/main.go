package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Passenger struct {
	PassengerID string `json:"PassengerID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	MobileNo    string `json:"MobileNo"`
	Email       string `json:"Email"`
}

type Driver struct {
	DriverID  string `json:"DriverID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	MobileNo  string `json:"MobileNo"`
	Email     string `json:"Email"`
	LicenseNo string `json:"LicenseNo"`
	Status    string `json:"Status"`
}

type Location struct {
	PickUpPostalCode  string `json:"PickUpPostalCode"`
	DropOffPostalCode string `json:"DropOffPostalCode"`
}

type Trip struct {
	TripID            string `json:"TripID"`
	PassengerID       string `json:"PassengerID"`
	DriverID          string `json:"DriverID"`
	PickUpPostalCode  string `json:"PickUpPostalCode"`
	DropOffPostalCode string `json:"DropOffPostalCode"`
	StartDateTime     string `json:"StartDateTime"`
	EndDateTime       string `json:"EndDateTime"`
}

func GetSingleRecord(db *sql.DB, tripID string) Trip {
	var foundTrip Trip
	query := fmt.Sprintf("Select * FROM DRide.Trip WHERE TripID = " + "'" + tripID + "'")

	err := db.QueryRow(query).Scan(&foundTrip.TripID, &foundTrip.PassengerID,
		&foundTrip.DriverID, &foundTrip.PickUpPostalCode, &foundTrip.DropOffPostalCode, &foundTrip.StartDateTime, &foundTrip.EndDateTime)

	if err != nil && err != sql.ErrNoRows {
		return foundTrip
	}

	return foundTrip
}

func trip(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRide")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)

	//Get Trip Record
	if r.Method == "GET" {

		retrievedTrip := GetSingleRecord(db, params["tripID"])
		json.NewEncoder(w).Encode(retrievedTrip)
		fmt.Println("Returned retrievd Driver")
		return
	}

	//Posting a request for trip
	if r.Header.Get("Content-type") == "application/json" {

		// Register
		if r.Method == "POST" {

			passengerID := params["ID"]
			var location Location
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &location)

				if location.PickUpPostalCode == "" || location.DropOffPostalCode == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Pick Up / Drop Off Location not entered"))

					defer db.Close()
					return
				}

				//Get available driver
				results, err := db.Query("Select * FROM Dride.Driver")

				if err != nil {
					panic(err.Error())
				}

				for results.Next() {
					// map this type to the record in the table
					var driver Driver
					err = results.Scan(&driver.DriverID, &driver.FirstName,
						&driver.LastName, &driver.MobileNo, &driver.Email, &driver.LicenseNo, &driver.Status)

					if err != nil {
						panic(err.Error())
					}
					if driver.Status == "Available" {
						query := fmt.Sprintf("INSERT INTO Trip (PassengerID, DriverID, PickUpPostalCode, DropOffPostalCode, StartDateTime) VALUES ('%s', '%s', '%s', '%s', '%s');",
							passengerID, driver.DriverID, location.PickUpPostalCode, location.DropOffPostalCode, time.Now().Format("2006-01-02 15:04:05"))

						_, err := db.Query(query)

						if err != nil {
							panic(err.Error())
						}

						fmt.Println("Trip has been created")
						// Setting driver availability to Unavailable
						queryStatus := fmt.Sprintf("UPDATE Driver SET Status = 'Unavailable' WHERE DriverID = '%s';",
							driver.DriverID)
						_, err2 := db.Query(queryStatus)

						if err2 != nil {
							panic(err2.Error())
						}

						fmt.Println("Changed driver status")
						w.WriteHeader(http.StatusCreated)
						json.NewEncoder(w).Encode(driver)
						//w.Write([]byte("Added as test"))
						return
					}

				}

				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("No available Drivers please try again later"))

			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver information " +
					"in JSON format"))
				defer db.Close()
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/trip/{ID}", trip).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 120")
	log.Fatal(http.ListenAndServe(":120", router))

	fmt.Println("Database opened")

}
