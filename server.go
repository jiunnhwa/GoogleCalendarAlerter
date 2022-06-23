package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func ServeRoutes() error {

	//VIEWS
	http.HandleFunc("/", home)

	fmt.Println("Started http listening...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		return err
	}
	return nil
}

//handles home page
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to My Google Calendar Alerter.")
	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RawQuery:", u.RawQuery)
	authCode := u.Query().Get("code")
	fmt.Println("home/authCode:", authCode)
	processTokenToFile(authCode)
}
