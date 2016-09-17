package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"model_1"
	"net/http"
	"time"
)

type (
	SensorDataCollectController struct {
		session *mgo.Session
		sqldb   *sql.DB
	}
)

func NewSensorDataCollectController(session *mgo.Session, sqldb *sql.DB) *SensorDataCollectController {
	return &SensorDataCollectController{session, sqldb}
}

func (sd SensorDataCollectController) SendSensorData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	go sendData(ps.ByName("sensorName"), ps.ByName("userName"), r, sd)
}

func sendData(sensorName string, userName string, req *http.Request, sd SensorDataCollectController) {

	sensorData := model.SensorData{}
	collectionName := userName + "-" + sensorName
	// JSON to Object Marshalling
	jsonDecoder := json.NewDecoder(req.Body)
	err := jsonDecoder.Decode(&sensorData)

	if err != nil {
		fmt.Errorf("Error in decoding the JSON: %v", err)
	}

	sensorData.Time = time.Now().UTC()
	sensorData.SensorId = bson.NewObjectId()
	sd.session.DB("mobisense").C(collectionName).Insert(sensorData)

}
