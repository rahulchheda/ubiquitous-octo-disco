package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

const (
	StartingPort = 12345
	PortRunning  = 12345
	OffSet       = 1
)

type Time struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type Day struct {
	Date     string `json:"date,omitempty"`
	Schedule []Time `json:"schedule,omitempty"`
}

var days []Day

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func GetDayEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range days {
		if item.Date == params["date"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Day{})
}

func GetDaysEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(days)
}

func CreateDaySchedule(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var day Day
	_ = json.NewDecoder(req.Body).Decode(&day)
	day.Date = params["date"]
	days = append(days, day)
	reqBody, err := json.Marshal(day)

	if err != nil {
		print(err)
	}
	ip := ReadUserIP(req)
	fmt.Println(ip)
	resp, err := http.Post("http://localhost:8080"+"/days/"+day.Date,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))
	json.NewEncoder(w).Encode(day)
}

func DeleteDaysEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range days {
		if item.Date == params["date"] {
			days = append(days[:index], days[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(days)
}

func main() {
	router := mux.NewRouter()
	days = append(days, Day{Date: "01122020", Schedule: []Time{{Start: 1100, End: 1200}, {Start: 1800, End: 1900}}})
	days = append(days, Day{Date: "01112020", Schedule: []Time{{Start: 1000, End: 1200}, {Start: 1600, End: 1900}}})
	router.HandleFunc("/days", GetDaysEndpoint).Methods("GET")
	router.HandleFunc("/days/{date}", GetDayEndpoint).Methods("GET")
	router.HandleFunc("/days/{date}", CreateDaySchedule).Methods("POST")
	router.HandleFunc("/days/{date}", DeleteDaysEndpoint).Methods("DELETE")
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(":"+fmt.Sprint(PortRunning), router))
	}()

	wg.Wait()
}
