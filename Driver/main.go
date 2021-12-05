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

type Driver struct {
	DriverID  string `json:"DriverID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	MobileNo  string `json:"MobileNo"`
	Email     string `json:"Email"`
	LicenseNo string `json:"LicenseNo"`
	Status    string `json:"Status"`
}

func validKey(r *http.Request, driverID string) bool {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRide")

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

func GetSingleRecord(db *sql.DB, driverID string) Driver {
	var foundDriver Driver
	query := fmt.Sprintf("Select * FROM DRide.Driver WHERE DriverID = " + "'" + driverID + "'")

	err := db.QueryRow(query).Scan(&foundDriver.DriverID, &foundDriver.FirstName,
		&foundDriver.LastName, &foundDriver.MobileNo, &foundDriver.Email, &foundDriver.LicenseNo, &foundDriver.Status)
	if err != nil && err != sql.ErrNoRows {
		return foundDriver
	}

	return foundDriver
}

func driver(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DRide")
	// handle db error
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)

	//Login
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

	if r.Header.Get("Content-type") == "application/json" {

		// Register
		if r.Method == "POST" {

			// read the string sent to the service
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
					// Add to database here
					query := fmt.Sprintf("INSERT INTO Driver VALUES ('%s', '%s', '%s', '%s', '%s', '%s', 'Available')",
						newDriver.DriverID, newDriver.FirstName, newDriver.LastName, newDriver.MobileNo, newDriver.Email, newDriver.LicenseNo)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Added as test"))

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

				// check if driverID exists; add only if
				// number does not exist
				retrivedDriver := GetSingleRecord(db, params["driverID"])

				if retrivedDriver.DriverID == "" {
					// Add to database here
					query := fmt.Sprintf("INSERT INTO Driver VALUES ('%s', '%s', '%s', '%s', '%s', '%s', 'Available')",
						newUpdatedDriverInfo.DriverID, newUpdatedDriverInfo.FirstName, newUpdatedDriverInfo.LastName, newUpdatedDriverInfo.MobileNo, newUpdatedDriverInfo.Email, newUpdatedDriverInfo.LicenseNo)

					_, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Added as test"))

					if err != nil {
						panic(err.Error())
					}
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver added: " +
						params["driverID"]))
				} else {
					// update course
					query := fmt.Sprintf("UPDATE Driver SET FirstName = '%s', LastName = '%s', MobileNo = '%s', Email = '%s', LicenseNumber = '%s' WHERE driverID = '%s'",
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
	router.HandleFunc("/api/driver/{driverID}", driver).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 100")
	log.Fatal(http.ListenAndServe(":100", router))

	fmt.Println("Database opened")

}
