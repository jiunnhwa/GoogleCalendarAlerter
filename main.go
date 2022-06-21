package main

import (
	"fmt"
	"time"

	"googlecal/service/gcal"
)

func init() {
	fmt.Println("*********************************************************")
	fmt.Println("Welcome to My Google Calendar Alerter.")
	fmt.Println("*********************************************************")
}
func main() {

	NOW := time.Now()
	startTime := NOW.Truncate(time.Minute).Add(time.Minute)
	duration := startTime.Sub(NOW)
	fmt.Println("Application start:", NOW)
	fmt.Println("Calendar Service starting at", startTime, " Wait seconds:", duration)
	time.Sleep(duration) //Sleep Until next minute

	myGCal := gcal.New("GCalService") //Create calendar service
	go func() {                       //On each start of minute:
		fmt.Printf("OnTickM1/StartTime: \t%v\n", startTime)
		ticker := time.NewTicker(time.Minute)
		for ; true; <-ticker.C {
			NOW = time.Now()
			fmt.Println("OnTickM1:", NOW)
			CURR := NOW.Format(time.RFC3339)
			TOM := NOW.AddDate(0, 0, 1).Format(time.RFC3339)
			myGCal.FetchCal(CURR, TOM)      //Fetch from Now till tomorrow
			myGCal.CheckForUpcomingEvents() //Check
		}
	}()

	ServeRoutes() //For OAuth callback

	fmt.Println("Press Enter To Exit")
	fmt.Scanln()

}
