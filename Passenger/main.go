package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

//Creation of Struct
type Passenger struct {
	PassengerID string
	FirstName   string
	LastName    string
	MobileNo    string
	Email       string
}

// Checking for correct login information
func validKey(r *http.Request, passengerID string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRidePassenger")

	if err != nil {
		panic(err.Error())
	}

	v := r.URL.Query()
	retrivedPassenger := GetSingleRecord(db, passengerID)
	if mobileNo, ok := v["mobileNo"]; ok {
		if mobileNo[0] == retrivedPassenger.MobileNo {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Getting of a single passenger record using passenger id
func GetSingleRecord(db *sql.DB, passengerID string) Passenger {
	var foundPassenger Passenger
	query := fmt.Sprintf("Select * FROM DRidePassenger.Passenger WHERE PassengerID = " + "'" + passengerID + "'")

	err := db.QueryRow(query).Scan(&foundPassenger.PassengerID, &foundPassenger.FirstName,
		&foundPassenger.LastName, &foundPassenger.MobileNo, &foundPassenger.Email)

	if err != nil && err != sql.ErrNoRows {
		return foundPassenger
	}

	return foundPassenger
}

// Passenger API with methods GET PUT POST DELETE
func passenger(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRidePassenger")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)

	// Login
	if r.Method == "GET" {
		if !validKey(r, params["passengerID"]) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("401 - Invalid Login information"))
			return
		}
		retrivedPassenger := GetSingleRecord(db, params["passengerID"])
		json.NewEncoder(w).Encode(retrivedPassenger)
		fmt.Println("Returned retrievd Passenger")
		return
	}

	// Delete
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - Unable to delete account due to auditing reasons"))
	}

	if r.Header.Get("Content-type") == "application/json" {

		// Register
		if r.Method == "POST" {

			// read the string sent to the service
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.PassengerID == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Passenger ID not entered"))

					defer db.Close()
					return
				}

				// check if passenger exists; add only if
				// passenger does not exist
				retrivedPassenger := GetSingleRecord(db, newPassenger.PassengerID)
				if retrivedPassenger.PassengerID == "" {
					// Insert new passenger details into database here
					query := fmt.Sprintf("INSERT INTO DRidePassenger.Passenger VALUES ('%s', '%s', '%s', '%s', '%s')",
						newPassenger.PassengerID, newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNo, newPassenger.Email)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Passenger Added"))

				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Passenger already exist"))
					defer db.Close()
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " +
					"in JSON format"))
				defer db.Close()
			}
		}

		// Updating passenger information
		if r.Method == "PUT" {
			// Ensures that the passenger is log-ed in
			if !validKey(r, params["passengerID"]) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("401 - Invalid key"))
				return
			}

			var newUpdatedPassengerInfo Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newUpdatedPassengerInfo)

				if params["passengerID"] == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply passenger " +
							" information " +
							"in JSON format"))
					return
				}

				// check if passenger exists; add only if
				// passenger does not exist
				retrivedPassenger := GetSingleRecord(db, params["passengerID"])

				if retrivedPassenger.PassengerID == "" {
					// Insert into database here
					query := fmt.Sprintf("INSERT INTO DRidePassenger.Passenger VALUES ('%s', '%s', '%s', '%s', '%s')",
						newUpdatedPassengerInfo.PassengerID, newUpdatedPassengerInfo.FirstName, newUpdatedPassengerInfo.LastName, newUpdatedPassengerInfo.MobileNo, newUpdatedPassengerInfo.Email)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Added passenger"))

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Passenger added: " +
						retrivedPassenger.PassengerID))
				} else {
					// update passenger information here
					query := fmt.Sprintf("UPDATE DRidePassenger.Passenger SET FirstName = '%s', LastName = '%s', MobileNo = '%s', Email = '%s' WHERE PassengerID = '%s'",
						newUpdatedPassengerInfo.FirstName, newUpdatedPassengerInfo.LastName, newUpdatedPassengerInfo.MobileNo, newUpdatedPassengerInfo.Email, params["passengerID"])

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("202 - Passenger updated: " +
						retrivedPassenger.PassengerID))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide passenger information"))
			}
		}

	}

	defer db.Close()
}

func main() {
	router := mux.NewRouter()
	// Passenger API Methods
	router.HandleFunc("/api/v1/passenger/{passengerID}", passenger).Methods(
		"GET", "PUT", "POST", "DELETE")

	// Using port 80 as passenger API
	fmt.Println("Listening at port 80")
	log.Fatal(http.ListenAndServe(":80", router))

	fmt.Println("Database opened")

}
