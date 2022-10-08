package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
)

const (
	glassPrice = 10
	icePrice   = 50
	StandPrice = 10
	days       = 7
	balance    = 2000
)

type UserInfo struct {
	Day        int64  `json:"day"`
	Balance    int64  `json:"balance"`
	Weather    string `json:"weather"`
	RainChance int64  `json:"rain_chance"`
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
	v, err := json.MarshalIndent(DayResponse{Users[id].Day}, "", "\t")
	_, err = w.Write(v)

	if err != nil {
		log.Println(err)
	}
}

func getWeather(w http.ResponseWriter, r *http.Request) {
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

	reasons := []string{
		"sunny",
		"hot",
		"cloudy",
	}
	n := rand.Int() % len(reasons)
	log.Println(reasons[n])

	Users[id] = UserInfo{Users[id].Day, Users[id].Balance, reasons[n], rand.Int63n(100)}
	v, err := json.MarshalIndent(WeatherResponce{reasons[n], Users[id].RainChance}, "", "\t")
	_, err = w.Write(v)

	if err != nil {
		log.Println(err)
	}
}

func setId(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	Users[id] = UserInfo{Day: 1, Balance: balance}
	w.WriteHeader(200)
	v, err := json.MarshalIndent(IdResponse{id}, "", "\t")
	_, err = w.Write(v)

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

	/*body, _ := io.ReadAll(r.Body)
	println(string(body))*/
	req := CalcRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err)
		return
	}

	users := int64(100)

	coef := 0.0

	if Users[id].Weather == "hot" {
		if req.IceAmount < 4 {
			users = users / (4 - req.IceAmount)
		} else {
			users = int64(float64(users) * 1.5)
		}
	}

	if Users[id].Weather == "sunny" {
		if req.IceAmount < 2 {
			users = users / (2 - req.IceAmount)
		} else if req.IceAmount <= 4 {
			users = int64(float64(users) * 1.5)
		} else {
			users = int64(float64(users) * 0.5)
		}
	}

	if Users[id].Weather == "cloudy" {
		if rand.Int63n(100) < Users[id].RainChance {
			if req.IceAmount == 0 {
				users = users
			} else {
				users = int64(float64(users) * 0.5)
			}
		} else {
			if req.IceAmount == 0 {
				users = int64(float64(users) / 1.5)
			} else if req.IceAmount <= 2 {
				users = int64(float64(users) * 1.5)
			} else {
				users = int64(float64(users) * 0.5)
			}
		}
	}

	if req.StandAmount > rand.Int63n(20*3)/2 {
		coef *= 2
	} else {
		coef /= 2
	}

	coef += rand.Float64()

	profit := int64(math.Min(float64(users), float64(req.CupsAmount))*float64(req.Price)) -
		glassPrice*req.CupsAmount - icePrice*req.IceAmount - StandPrice*req.StandAmount

	Users[id] = UserInfo{Users[id].Day + 1, Users[id].Balance + profit, "", 0}

	w.WriteHeader(200)
	v, err := json.MarshalIndent(CalcResponce{Users[id].Balance, profit, Users[id].Day}, "", "\t")
	_, err = w.Write(v)

	if err != nil {
		log.Println(err)
	}
	if Users[id].Day == days {
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
