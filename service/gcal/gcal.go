/*
	Package implements the base service methods for google calendar.
*/

package gcal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"googlecal/service"
)

//Job is a record to hold job details
type Job struct {
	Name               string
	StartTime, EndTime time.Time
	Request            string
	Response           string
	Result             string
	//Error              Status

	CalService  *calendar.Service
	CalEvents   *calendar.Events
	WebClient   *http.Client
	OauthConfig *oauth2.Config
	OauthToken  *oauth2.Token
}

const CREDENTIAL_FILE, TOKEN_FILE = "credentials.json", "token.json"
const ALERT_FILE_ON_M5, ALERT_FILE_ON_TIME = "end.mp3", "end.mp3"

//NewJob constructor
func New(name string) *Job {
	j := &Job{Name: name}

	//Parse CREDENTIAL_FILE
	bytes, err := ioutil.ReadFile(CREDENTIAL_FILE)
	if err != nil {
		log.Fatalf("Unable to read credential file: %j", err)
	}
	j.OauthConfig, err = google.ConfigFromJSON(bytes, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse credential file to config: %v", err)
	}

	//Get token from file or web,
	if err := j.getTokenFromFile(TOKEN_FILE); err != nil {
		j.getTokenFromWeb().saveToken(TOKEN_FILE)
	}

	//create client
	j.WebClient = j.OauthConfig.Client(context.Background(), j.OauthToken)

	//create calendar service
	j.CalService, err = calendar.NewService(context.Background(), option.WithHTTPClient(j.WebClient))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return j
}

// Request a token from the web, then returns the retrieved token.
func (j *Job) getTokenFromWeb() *Job {
	authURL := j.OauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Enter the authorization code returned from the authorization line: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := j.OauthConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	j.OauthToken = tok
	return j
}

// Retrieves a token from a local file.
func (j *Job) getTokenFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Unable to open file: %v", err)
		return err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		log.Fatalf("Unable to decode: %v", err)
		return err
	}
	j.OauthToken = tok
	return nil
}

// Saves a token to a file path.
func (j *Job) saveToken(path string) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(j.OauthToken)
}

//FetchCal fetches the primary calendar with time from  and till
func (j *Job) FetchCal(from, till string) {
	events, err := j.CalService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(from).TimeMax(till).MaxResults(100).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	//fmt.Println(events)
	j.CalEvents = events
}

//CheckForUpcomingEvents prints and play sound upon event reached.
func (j *Job) CheckForUpcomingEvents() {
	fmt.Println("Upcoming events:")
	if len(j.CalEvents.Items) == 0 {
		fmt.Println("There are no upcoming events from now till tomorrow...")
	} else {
		for _, item := range j.CalEvents.Items {
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
				fmt.Println("M10:", t2.Format("2006-01-02 15:04:05"))
				service.PlaySound(ALERT_FILE_ON_M5)
			}

			if int(timediff.Minutes()) == 0 {
				fmt.Println("NOW:", t2.Format("2006-01-02 15:04:05"))
				service.PlaySound(ALERT_FILE_ON_TIME)
			}
		}
	}
}
