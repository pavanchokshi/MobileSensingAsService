package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

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

	fmt.Println("Sensor Server Moldule using port:8088")
	router := httprouter.New()
	router.Handle("POST", "/sendsensordata/:userName/:sensorName", SendSensorData)
	log.Fatal(http.ListenAndServe(":8088", router))
}

func getDbSession(collectionName string) *mgo.Session {

	info := mgo.CollectionInfo{}
	info.Capped = false

	// Connect to our mongoDB
	s, err := mgo.Dial("mongodb://user01:user01@ds031347.mongolab.com:31347/mobisense_phy")

	info.Capped = true
	info.MaxBytes = 16
	info.MaxDocs = 1

	collectionName += "_" + "latest"
	s.DB("mobisense_phy").C(collectionName).Create(&info)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

func SendSensorData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sensorData := SensorData{}

	// JSON to Object Marshalling
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&sensorData)

	if err != nil {
		fmt.Errorf("Error in decoding the JSON: %v", err)
	}

	sensorData.Time = time.Now().UTC()
	sensorData.SensorId = bson.NewObjectId()

	go sendData(ps.ByName("sensorName"), ps.ByName("userName"), sensorData)
}

func sendData(sensorName string, userName string, sensorData SensorData) {

	collectionName := sensorName
	session := getDbSession(collectionName)

	fmt.Println(sensorData)
	session.DB("mobisense_phy").C(collectionName).Insert(sensorData)
	collectionName += "_" + "latest"
	session.DB("mobisense_phy").C(collectionName).Insert(sensorData)
	session.Close()
}
