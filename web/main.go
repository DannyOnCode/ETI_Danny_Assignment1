package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

//Page 1 - Main Page
func mainMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("mainMenu.html"))
	tmpl.Execute(w, struct{ Success bool }{true})
}

// Login page maybe
func loginDriver(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login/driverLogin.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := ContactDetails{
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	// TODO: Connect to API in Driver Microservice (Look into useRest Tutorial)
	fmt.Println(details)
	tmpl.Execute(w, details)

}

func loginPassenger(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login/passengerLogin.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := ContactDetails{
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	// TODO: Connect to API in Passenger Microservice (Look into useRest Tutorial)
	fmt.Println(details)
	tmpl.Execute(w, details)

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", mainMenu)
	r.HandleFunc("/login/driver", loginDriver)
	r.HandleFunc("/login/passenger", loginPassenger)
	fmt.Println("Listening at port 5000")
	http.ListenAndServe(":5000", r)
}
