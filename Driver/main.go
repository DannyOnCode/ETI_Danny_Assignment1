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
type Driver struct {
	DriverID  string `json:"DriverID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	MobileNo  string `json:"MobileNo"`
	Email     string `json:"Email"`
	LicenseNo string `json:"LicenseNo"`
	Status    string `json:"Status"`
}

// Checking for correct login information
func validKey(r *http.Request, driverID string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRideDriver")

	if err != nil {
		panic(err.Error())
	}

	v := r.URL.Query()
	retrivedDriver := GetSingleRecord(db, driverID)
	if mobileNo, ok := v["mobileNo"]; ok {
		if mobileNo[0] == retrivedDriver.MobileNo {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Getting a single driver record using driver id
func GetSingleRecord(db *sql.DB, driverID string) Driver {
	var foundDriver Driver
	query := fmt.Sprintf("Select * FROM DRideDriver.Driver WHERE DriverID = " + "'" + driverID + "'")

	err := db.QueryRow(query).Scan(&foundDriver.DriverID, &foundDriver.FirstName,
		&foundDriver.LastName, &foundDriver.MobileNo, &foundDriver.Email, &foundDriver.LicenseNo, &foundDriver.Status)
	if err != nil && err != sql.ErrNoRows {
		return foundDriver
	}

	return foundDriver
}

// Driver API with methods GET PUT POST DELETE
func driver(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRideDriver")
	// handle db error
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)

	// Login
	if r.Method == "GET" {
		if !validKey(r, params["driverID"]) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("401 - Invalid Login information"))
			return
		}
		retrivedDriver := GetSingleRecord(db, params["driverID"])
		json.NewEncoder(w).Encode(retrivedDriver)
		fmt.Println("Returned retrievd Driver")
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

			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.DriverID == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Driver ID not entered"))

					defer db.Close()
					return
				}

				// check if driver exists; add only if
				// driver does not exist
				retrivedDriver := GetSingleRecord(db, newDriver.DriverID)
				if retrivedDriver.DriverID == "" {
					// Insert driver into database here
					query := fmt.Sprintf("INSERT INTO DRideDriver.Driver VALUES ('%s', '%s', '%s', '%s', '%s', '%s', 'Available')",
						newDriver.DriverID, newDriver.FirstName, newDriver.LastName, newDriver.MobileNo, newDriver.Email, newDriver.LicenseNo)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Driver Added"))

				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Driver already exist"))
					defer db.Close()
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver information " +
					"in JSON format"))
				defer db.Close()
			}
		}

		// Updating Driver information
		if r.Method == "PUT" {

			if !validKey(r, params["driverID"]) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("401 - Invalid key"))
				return
			}

			var newUpdatedDriverInfo Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newUpdatedDriverInfo)

				if params["driverID"] == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply Driver " +
							" information " +
							"in JSON format"))
					return
				}

				// check if driver exists; add only if
				// driver does not exist
				retrivedDriver := GetSingleRecord(db, params["driverID"])

				if retrivedDriver.DriverID == "" {
					// Insert driver into database here
					query := fmt.Sprintf("INSERT INTO DRideDriver.Driver VALUES ('%s', '%s', '%s', '%s', '%s', '%s', 'Available')",
						newUpdatedDriverInfo.DriverID, newUpdatedDriverInfo.FirstName, newUpdatedDriverInfo.LastName, newUpdatedDriverInfo.MobileNo, newUpdatedDriverInfo.Email, newUpdatedDriverInfo.LicenseNo)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Driver Added"))

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver added: " +
						params["driverID"]))
				} else {
					// update driver information here
					query := fmt.Sprintf("UPDATE DRideDriver.Driver SET FirstName = '%s', LastName = '%s', MobileNo = '%s', Email = '%s', LicenseNumber = '%s' WHERE driverID = '%s'",
						newUpdatedDriverInfo.FirstName, newUpdatedDriverInfo.LastName, newUpdatedDriverInfo.MobileNo, newUpdatedDriverInfo.Email, newUpdatedDriverInfo.LicenseNo, params["driverID"])

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("202 - Driver updated: " +
						params["driverID"]))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please provide Driver information"))
			}
		}

	}

	defer db.Close()
}

func main() {
	router := mux.NewRouter()
	// Driver API methods
	router.HandleFunc("/api/v1/driver/{driverID}", driver).Methods(
		"GET", "PUT", "POST", "DELETE")

	// Using port 100 as Driver API
	fmt.Println("Listening at port 100")
	log.Fatal(http.ListenAndServe(":100", router))

	fmt.Println("Database opened")

}
