package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	Sensor struct {
		SensorId     string    `json:"sensorid,omitempty"`
		OwnerId      string    `json:"ownerid"`
		SensorName   string    `json:"sensorname"`
		SensorDesc   string    `json:"sensordesc"`
		SensorType   string    `json:"sensortype"`
		IsShared     bool      `json:"isshared"`
		IsEnabled    bool      `json:"isenabled,omitempty"`
		IsActivated  bool      `json:"isactivated,omitempty"`
		RegisterTime time.Time `json:"created,omitempty"`
		Latitude     string    `json:"latitude,omitempty"`
		Longitude    string    `json:"longitude,omitempty"`
	}

	SensorAuthenticate struct {
		OwnerId    string `json:"ownerid"`
		SensorName string `json:"sensorname"`
		Password   string `json:"password"`
	}

	SensorData struct {
		SensorId   bson.ObjectId `json:"id" bson:"_id,omitempty"`
		SensorName string        `json:"sensorname"`
		SensorData string        `json:"sensordata"`
		Location   struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		Time time.Time `json:"created,omitempty"`
	}

	VirtSensor struct {
		VirtSensorId       string    `json:"virtsensorid,omitempty"`
		VirtSensorName     string    `json:"virtsensorname"`
		VirtSensorUserId   string    `json:"virtsensoruserid"`
		PhysicalSensorName string    `json:"sensorname"`
		ProjectId	   string    `json:"projectid"`
		IsEnabled          bool      `json:"isenabled,omitempty"`
		DataFreqinSec      int       `json:"datafreqinsec"`
		UsageCounterinMins int       `json:"usagecounterinmins,omitempty"`
		RegisterTime       time.Time `json:"created,omitempty"`
		Status             int       `json:"status"`
	}
)
