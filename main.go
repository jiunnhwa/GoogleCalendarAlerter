package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"googlecal/service"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const MAX_TRIES = 5
const MAX_RESULTS = 100
const CREDENTIAL_FILE, TOKEN_FILE = "credentials.json", "token.json"
const ALERT_FILE_ON_M5, ALERT_FILE_ON_TIME = "end.mp3", "end.mp3"

//custom errors
var (
	ErrFetchCalendarEvents      = errors.New("fail to retrieve calendar's events")
	ErrGetTokenFromWeb          = errors.New("fail to get token from web")
	ErrSaveTokenToFile          = errors.New("fail to save to token file")
	ErrReadingCredenialFile     = errors.New("fail to read credential file")
	ErrParsingCredenialFile     = errors.New("fail to parse credential file")
	ErrGetCalendarClient        = errors.New("fail to get calendar client")
	ErrTimeOutWaitingForService = errors.New("time out waiting for service...")
	ErrTokenExpired             = errors.New("oauth2: token expired and refresh token is not set")
	//
)

var myConfig *oauth2.Config
var myClient *http.Client
var myGCal *calendar.Service
var myEvents *calendar.Events
var myToken *oauth2.Token

func init() {
	fmt.Println("*********************************************************")
	fmt.Println("Welcome to My Google Calendar Alerter.")
	fmt.Println("*********************************************************")

	log.Println("Application start:", time.Now())
}

func main() {
	go ServeRoutes() //For OAuth callback
	getCredFromFile()
START:
	if _, err := os.Stat(TOKEN_FILE); err == nil { // TOKEN_FILE exists
		getClientService()
	} else {
		getTokenFromWeb()
		goto START
	}

	//tries 5 times, with increasing backoff, await service to be ready
	maxTries := 1
	for ; maxTries <= MAX_TRIES; maxTries++ {
		time.Sleep(10 * time.Duration(maxTries) * time.Second) // increasing backoff
		fmt.Println("maxTries:", maxTries, time.Now())
		if myGCal != nil {
			go FetchChecker()
			break
		}
	}
	if maxTries >= MAX_TRIES {
		fmt.Println("maxTries:", maxTries, ErrTimeOutWaitingForService.Error())
	}

	fmt.Println("Press Enter To Exit")
	fmt.Scanln()
}

//get calendar service from existing token
func getClientService() {
	myClient = myConfig.Client(context.Background(), getExistingTokenFromFile())
	svc, err := calendar.NewService(context.Background(), option.WithHTTPClient(myClient))
	if err != nil {
		log.Fatalf("%s: %v", ErrGetCalendarClient.Error(), err)
	} else {
		myGCal = svc
		fmt.Println("myGCal:", myGCal.BasePath)
	}
}

//read existing credential
func getCredFromFile() {
	if bytes, err := readFile(CREDENTIAL_FILE); err != nil {
		log.Fatalf("%s: %v", ErrReadingCredenialFile.Error(), err)
	} else {
		myConfig, err = google.ConfigFromJSON(bytes, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("%s: %v", ErrParsingCredenialFile.Error(), err)
		}
	}
}

//read and unmarshal to token
func getExistingTokenFromFile() *oauth2.Token {
	token := &oauth2.Token{}
	if bytes, err := readFile(TOKEN_FILE); err == nil {
		if err := json.Unmarshal(bytes, &token); err != nil {
			panic(err)
		}
	}
	return token
}

//get token from url
func getTokenFromWeb() {
	authURL := myConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("authorization url: \n%v\n", authURL)
	//exec.Command("cmd", "/C", "start", authURL).Run() //Authorization Error	Error 400: invalid_request	Required parameter is missing: response_type
	webbrowser.Open(authURL)

	//await sign-in, exchange token, ...
	for {
		if _, err := os.Stat(TOKEN_FILE); err == nil { // TOKEN_FILE is renewed, exists
			fmt.Println("TOKEN_FILE downloaded.")
			return
		} else {
			fmt.Println("TOKEN_FILE missing.")
		}
		time.Sleep(time.Second) //wait
	}
}

//read the given file and return it as bytes and error
func readFile(fname string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}
	return (bytes), nil
}

//fetch and do alert
func FetchChecker() { //On each start of minute:
	now := time.Now()
	startTime := now.Truncate(time.Minute).Add(time.Minute)
	duration := startTime.Sub(now)
	fmt.Println("FetchChecker Service starting at", startTime, " Wait seconds:", duration)
	time.Sleep(duration) //Sleep Until next minute

	fmt.Printf("FetchChecker:StartTime: \t%v\n", now)
	ticker := time.NewTicker(time.Minute)
	for ; true; <-ticker.C {
		now = time.Now()
		fmt.Println("OnTickM1:", now)
		CURR := now.Format(time.RFC3339)
		TOM := now.AddDate(0, 0, 1).Format(time.RFC3339)
		FetchCal(CURR, TOM)      //Fetch from Now till tomorrow
		CheckForUpcomingEvents() //Check
	}
}

//FetchCal fetches the primary calendar with time from and till
func FetchCal(from, till string) {
	events, err := myGCal.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(from).TimeMax(till).MaxResults(MAX_RESULTS).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("%s: %v", ErrFetchCalendarEvents.Error(), err)
	}
	myEvents = events
}

//CheckForUpcomingEvents prints and play sound upon event reached.
func CheckForUpcomingEvents() {
	fmt.Println("Upcoming events:")
	if len(myEvents.Items) == 0 {
		fmt.Println("There are no upcoming events from now till tomorrow...")
	} else {
		for _, item := range myEvents.Items {
			date := item.Start.DateTime
			fmt.Println("======================================")
			fmt.Printf("%v\n %v\n (%v)\n (%v)\n", service.ObfuscateLast4Char(item.Id), item.Summary, service.ObfuscateLast4Char(item.HtmlLink), item.Start)

			t2, _ := time.Parse(time.RFC3339, date)
			timediff := time.Since(t2)
			if timediff > 0 {
				fmt.Println("over since ", int(timediff.Minutes()), "mins ago. Event Time:", t2.Format("2006-01-02 15:04:05"))
			} else {
				fmt.Println("upcoming in", int(math.Abs(timediff.Minutes())), "mins. Event Time:", t2.Format("2006-01-02 15:04:05"))
			}
			fmt.Println("")

			if int(timediff.Minutes()) == -5 {
				fmt.Println("M15:", t2.Format("2006-01-02 15:04:05"))
				service.PlaySound(ALERT_FILE_ON_M5)
			}

			if int(timediff.Minutes()) == 0 {
				fmt.Println("NOW:", t2.Format("2006-01-02 15:04:05"))
				service.PlaySound(ALERT_FILE_ON_TIME)
			}
		}
	}
}

//exchange and save token
func processTokenToFile(authCode string) {
	//may have double callback with 2nd one empty authCode, and causing error oauth2: cannot fetch token: 400 Bad Request "error": "invalid_request "error_description": "Missing required parameter: code"
	fmt.Println("processTokenToFile:authCode:", authCode, time.Now())
	//Exchange for token
	if len(authCode) > 0 {
		token, err := myConfig.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Printf("%s: %v", ErrGetTokenFromWeb.Error(), err) //ignore re-entrant with empty authcode
		}
		myClient = myConfig.Client(context.Background(), token)
		saveTokenToFile(token)
	}
}

//Save the refresh token
func saveTokenToFile(token *oauth2.Token) {
	fmt.Println("Saving refresh token:", token)
	f, err := os.OpenFile(TOKEN_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("%s: %v", ErrSaveTokenToFile.Error(), err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
