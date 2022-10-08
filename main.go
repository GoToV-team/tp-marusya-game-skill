package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type UserInfo struct {
	Day     int64 `json:"day"`
	Balance int64 `json:"balance"`
}

type DayResponse struct {
	Day int64 `json:"day"`
}

type IdResponse struct {
	Id string `json:"id"`
}

type CalcRequest struct {
	CupsAmount  int64 `json:"cups_amount"`
	IceAmount   int64 `json:"ice_amount"`
	StandAmount int64 `json:"stand_amount"`
	Price       int64 `json:"price"`
}

type WeatherResponce struct {
	Weather string `json:"weather_name"`
	Chance  int64  `json:"rain_chance"`
}

type CalcResponce struct {
	Balance int64 `json:"balance"`
	Profit  int64 `json:"profit"`
	Day     int64 `json:"day"`
}

var Users = make(map[string]UserInfo)

func getDay(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(400)
		log.Println("getDay 400")
		return
	}

	if _, ok := Users[id]; !ok {
		w.WriteHeader(404)
		log.Println("getDay 404")
		return
	}

	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(DayResponse{Users[id].Day})

	if err != nil {
		log.Println(err)
	}
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	reasons := []string{
		"sunny",
		"hot",
		"cloudy",
	}
	n := rand.Int() % len(reasons)
	fmt.Print("Gonna work from home...", reasons[n])
	err := json.NewEncoder(w).Encode(WeatherResponce{reasons[n], rand.Int63n(100)})

	if err != nil {
		log.Println(err)
	}
}

func setId(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	Users[id] = UserInfo{Day: 1, Balance: 100}
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(IdResponse{id})

	if err != nil {
		log.Println(err)
	}
}

func calcValue(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(400)
		log.Println("getDay 400")
		return
	}

	if _, ok := Users[id]; !ok {
		w.WriteHeader(404)
		log.Println("getDay 404")
		return
	}

	req := CalcRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	profit := rand.Float64() * float64(req.CupsAmount) * float64(req.Price)

	profit *= rand.Float64()

	Users[id] = UserInfo{Users[id].Day + 1, Users[id].Balance + int64(profit)}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(CalcResponce{Users[id].Balance, int64(profit), Users[id].Day})

	if err != nil {
		log.Println(err)
	}
	if Users[id].Day == 7 {
		delete(Users, id)
	}
}

func main() {
	http.HandleFunc("/day", getDay)
	http.HandleFunc("/weather", getWeather)
	http.HandleFunc("/id", setId)
	http.HandleFunc("/calculate", calcValue)

	err := http.ListenAndServe(":80", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
