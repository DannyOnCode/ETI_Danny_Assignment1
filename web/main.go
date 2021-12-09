package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const passengerURL = "http://localhost:80/api/passenger"
const driverURL = "http://localhost:100/api/driver"
const tripURL = "http://localhost:120/api/trip"

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

var currentPassengerInfo Passenger
var currentDriverInfo Driver

//Page 1 - Main Page
func mainMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("mainMenu.html"))
	tmpl.Execute(w, struct{ Success bool }{true})
}

func pMainMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("passengerMainPage.html"))
	tmpl.Execute(w, currentPassengerInfo)
}

func dMainMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("driverMainPage.html"))
	tmpl.Execute(w, currentDriverInfo)
}

// Login Driver maybe
func loginDriver(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login/driverLogin.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Driver{
		DriverID: r.FormValue("driverid"),
		MobileNo: r.FormValue("contact"),
	}

	var url string
	if details.MobileNo != "" && details.DriverID != "" {
		url = driverURL + "/" + details.DriverID + "?mobileNo=" + details.MobileNo
	}

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		json.Unmarshal(data, &currentDriverInfo)
		response.Body.Close()
	}

	if currentDriverInfo.MobileNo != "" {
		http.Redirect(w, r, "http://localhost:5000/driver/main", http.StatusFound)
	}
	tmpl.Execute(w, currentDriverInfo)
}

//Register Driver
func registerDriver(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("register/driverRegister.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Driver{
		DriverID:  r.FormValue("driverid"),
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
		MobileNo:  r.FormValue("contact"),
		Email:     r.FormValue("email"),
		LicenseNo: r.FormValue("licenseno"),
	}

	// TODO: Connect to API in Passenger Microservice (Look into useRest Tutorial)
	jsonValue, _ := json.Marshal(details)

	response, err := http.Post(driverURL+"/"+details.MobileNo,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	fmt.Println(details)
	tmpl.Execute(w, details)

}

// Login Passenger page
func loginPassenger(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login/passengerLogin.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Passenger{
		PassengerID: r.FormValue("passengerid"),
		MobileNo:    r.FormValue("contact"),
	}

	var url string
	if details.MobileNo != "" && details.PassengerID != "" {
		url = passengerURL + "/" + details.PassengerID + "?mobileNo=" + details.MobileNo
	}
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		json.Unmarshal(data, &currentPassengerInfo)
		response.Body.Close()
	}

	if currentPassengerInfo.MobileNo != "" {
		http.Redirect(w, r, "http://localhost:5000/passenger/main", http.StatusFound)
	}
	tmpl.Execute(w, currentPassengerInfo)
}

func registerPassenger(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("register/passengerRegister.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Passenger{
		PassengerID: r.FormValue("passengerid"),
		FirstName:   r.FormValue("firstname"),
		LastName:    r.FormValue("lastname"),
		MobileNo:    r.FormValue("contact"),
		Email:       r.FormValue("email"),
	}

	// TODO: Connect to API in Passenger Microservice (Look into useRest Tutorial)
	jsonValue, _ := json.Marshal(details)

	response, err := http.Post(passengerURL+"/"+details.MobileNo,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	fmt.Println(details)
	tmpl.Execute(w, details)

}

func updatePassenger(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("updateInformation/passengerUpdate.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Passenger{
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
		MobileNo:  r.FormValue("contact"),
		Email:     r.FormValue("email"),
	}

	// TODO: Connect to API in Passenger Microservice (Look into useRest Tutorial)
	jsonValue, _ := json.Marshal(details)

	request, err := http.NewRequest(http.MethodPut, passengerURL+"/"+currentPassengerInfo.PassengerID+"?mobileNo="+currentPassengerInfo.MobileNo, bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	fmt.Println(details)
	tmpl.Execute(w, details)

}

func updateDriver(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("updateInformation/driverUpdate.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Driver{
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
		MobileNo:  r.FormValue("contact"),
		Email:     r.FormValue("email"),
		LicenseNo: r.FormValue("licenseno"),
	}

	// TODO: Connect to API in Passenger Microservice (Look into useRest Tutorial)
	jsonValue, _ := json.Marshal(details)

	request, err := http.NewRequest(http.MethodPut, driverURL+"/"+currentDriverInfo.DriverID+"?mobileNo="+currentDriverInfo.MobileNo, bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

	fmt.Println(details)
	tmpl.Execute(w, details)

}

func deletePassenger(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest(http.MethodDelete,
		passengerURL+"/"+currentPassengerInfo.PassengerID, nil)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	http.Redirect(w, r, "http://localhost:5000/passenger/main", http.StatusFound)
}

func deleteDriver(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest(http.MethodDelete,
		driverURL+"/"+currentDriverInfo.DriverID, nil)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	http.Redirect(w, r, "http://localhost:5000/driver/main", http.StatusFound)
}

func tripPassenger(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("Trip/passengerTrip.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Location{
		PickUpPostalCode:  r.FormValue("pickup"),
		DropOffPostalCode: r.FormValue("dropoff"),
	}

	var tripDriver Driver

	jsonValue, _ := json.Marshal(details)
	response, err := http.Post(tripURL+"/"+currentPassengerInfo.PassengerID,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		json.Unmarshal(data, &tripDriver)
		response.Body.Close()

	}

	tmpl.Execute(w, tripDriver)
}

func tripDriver(w http.ResponseWriter, r *http.Request) {
	var availableTrip Trip
	if r.Method != http.MethodPost {
		//TO DO: Add Get Function here to check if driver has ride available
		//Add check here if has ride, display ride. Else : template.Must(template.ParseFiles("Trip/driverNoTrip.html"))
		var url string
		if currentDriverInfo.DriverID != "" {
			url = tripURL + "/" + currentDriverInfo.DriverID + "?userType=" + "driver"
		} else {
			// Redirect to No trip page
			template.Must(template.ParseFiles("login/driverLogin.html"))
		}
		response, err := http.Get(url)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(response.StatusCode)
			fmt.Println(string(data))

			json.Unmarshal(data, &availableTrip)
			fmt.Println(availableTrip)
			response.Body.Close()
		}
		if availableTrip.TripID == "" {
			tmpl := template.Must(template.ParseFiles("Trip/driverNoTrip.html"))
			tmpl.Execute(w, nil)
		} else {
			tmpl := template.Must(template.ParseFiles("Trip/driverTrip.html"))
			tmpl.Execute(w, availableTrip)
		}
		return
	}
	// TO DO: Add PUT request here to start ride
	var url string
	if currentDriverInfo.DriverID != "" {
		url = tripURL + "/" + currentDriverInfo.DriverID + "?userType=" + "driver"
	} else {
		// Redirect to No trip page
		template.Must(template.ParseFiles("login/driverLogin.html"))
	}
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		json.Unmarshal(data, &availableTrip)
		fmt.Println(availableTrip)
		response.Body.Close()
	}
	jsonValue, _ := json.Marshal(availableTrip)
	request, err := http.NewRequest(http.MethodPut, tripURL+"/"+currentDriverInfo.DriverID, bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err = client.Do(request)

	var retrivedTrip Trip
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		json.Unmarshal(data, &retrivedTrip)
		response.Body.Close()
	}
	tmpl := template.Must(template.ParseFiles("Trip/startTripPage.html"))
	fmt.Println(retrivedTrip)
	tmpl.Execute(w, retrivedTrip)
}

func viewPassengerHistory(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("Trip/passengerHistory.html"))
	var url string
	if currentPassengerInfo.PassengerID != "" {
		url = tripURL + "/" + currentPassengerInfo.PassengerID + "?userType=" + "passenger"
	}
	response, err := http.Get(url)

	var tripArray []Trip
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		json.Unmarshal(data, &tripArray)
		response.Body.Close()
	}

	tmpl.Execute(w, tripArray)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", mainMenu)
	r.HandleFunc("/login/driver", loginDriver)
	r.HandleFunc("/login/passenger", loginPassenger)
	r.HandleFunc("/register/passenger", registerPassenger)
	r.HandleFunc("/register/driver", registerDriver)
	r.HandleFunc("/passenger/main", pMainMenu)
	r.HandleFunc("/driver/main", dMainMenu)
	r.HandleFunc("/update/passenger", updatePassenger)
	r.HandleFunc("/update/driver", updateDriver)
	r.HandleFunc("/delete/passenger", deletePassenger)
	r.HandleFunc("/delete/driver", deleteDriver)
	r.HandleFunc("/Trip/passengerTrip", tripPassenger)
	r.HandleFunc("/Trip/driverTrip", tripDriver)
	r.HandleFunc("/passenger/viewHistory", viewPassengerHistory)
	fmt.Println("Listening at port 5000")
	http.ListenAndServe(":5000", r)
}
