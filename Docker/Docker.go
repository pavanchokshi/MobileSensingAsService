package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

type Configuration struct {
	NumberOfSensors int      `json:"numsensor"`
	SensorId        []string `json:"sensorid"`
	Status          []bool   `json:"status"`
	PhySensorName   []string `json:"physensorname"`
	Frequency       []int    `json:"freq"`
}

var config Configuration
var session *mgo.Session

type SensorData struct {
	SensorId   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	SensorData string        `json:"sensordata"`
	Location   struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Time time.Time `json:"created,omitempty"`
}

func main() {

	fmt.Println("Starting Docker..")
	router := httprouter.New()
	router.Handle("PUT", "/startsensor/:sensorId", StartSensor)
	router.Handle("PUT", "/stopsensor/:sensorId", StopSensor)
	router.Handle("PUT", "/updatefrequency/:sensorId/:freq", UpdateFrequency)
	router.Handle("GET", "/getsensingdata/:sensorId", GetSensingData)
	router.Handle("GET", "/isalive", IsAlive)
	router.Handle("POST", "/getconfiguration", GetConfiguration)
	http.ListenAndServe(":9999", router)
}

func UpdateFrequency(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	for k, v := range config.SensorId {
		if v == ps.ByName("sensorId") {
			config.Frequency[k], _ = strconv.Atoi(ps.ByName("freq"))
		}
	}

}

func GetSensingData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sensorData := SensorData{}
	collectionName := ps.ByName("sensorId") + "_" + "latest"
	session, _ := mgo.Dial("mongodb://user01:user01@ds051170.mongolab.com:51170/mobisense")
	session.DB("mobisense").C(collectionName).Find(nil).One(&sensorData)
	session.Close()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sensorData)
}

func getDbSession() *mgo.Session {

	info := mgo.CollectionInfo{}
	info.Capped = false

	// Connect to our mongoDB
	s, err := mgo.Dial("mongodb://user01:user01@ds051170.mongolab.com:51170/mobisense")
	i := 0

	for i < config.NumberOfSensors {
		s.DB("mobisense").C(config.SensorId[i]).Create(&info)
		i++
	}

	info.Capped = true
	info.MaxBytes = 16
	info.MaxDocs = 1

	var collectionName string

	i = 0
	for i < config.NumberOfSensors {
		collectionName = config.SensorId[i] + "_" + "latest"
		s.DB("mobisense").C(collectionName).Create(&info)
		i++
	}

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}

func setSensorStatus(id string, val bool) int {

	for k, v := range config.SensorId {
		if v == id {
			config.Status[k] = val
			return k
		}
	}
	return -1
}

func collectData(index int) {
	session_phy, _ := mgo.Dial("mongodb://user01:user01@ds031347.mongolab.com:31347/mobisense_phy")
	session, _ := mgo.Dial("mongodb://user01:user01@ds051170.mongolab.com:51170/mobisense")
	for config.Status[index] == true {

		time.Sleep(time.Duration(config.Frequency[index]) * 1000 * time.Millisecond)
		collectionName := config.PhySensorName[index] + "_latest"
		sensorData := SensorData{}
		session_phy.DB("mobisense_phy").C(collectionName).Find(nil).One(&sensorData)

		session.DB("mobisense").C(fmt.Sprintf("%s_latest", config.SensorId[index])).Insert(sensorData)
		session.DB("mobisense").C(config.SensorId[index]).Insert(sensorData)
		fmt.Println(sensorData)
	}
	session_phy.Close()
	session.Close()
}

func IsAlive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func GetConfiguration(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&config)
	if err != nil {
		fmt.Errorf("Error in decoding the JSON: %v", err)
	}
	session = getDbSession()
}

func StartSensor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	SensorID := ps.ByName("sensorId")
	index := setSensorStatus(SensorID, true)
	if index != -1 {
		go collectData(index)
	}

}

func StopSensor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	SensorID := ps.ByName("sensorId")
	_ = setSensorStatus(SensorID, false)
}
