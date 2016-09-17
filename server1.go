package main

import (
	"controller_1"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecheney/profile"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	Database struct {
		Dbtype   string `json:"dbtype"`
		Dbname   string `json:"dbname"`
		Password string `json:"password"`
		User     string `json:"user"`
	} `json:"database"`
}

var db *sql.DB

func getDbSession() *sql.DB {

	// Parsing the Configuration File
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error decoding configuration:", err)
	}

	// Connecting to DB
	dbConnectionUrl := configuration.Database.User + ":" + configuration.Database.Password + "@/" + configuration.Database.Dbname + "?parseTime=true"
	db, err = sql.Open(configuration.Database.Dbtype, dbConnectionUrl)
	if err != nil {
		fmt.Errorf("Error in connecting to database: %v", err)
	}
	db.SetMaxIdleConns(50)
	return db

}

func main() {

	fmt.Println("Server Starting on port:8080")
	defer profile.Start(profile.CPUProfile).Stop()
	pc := controller.NewProjectController(getDbSession())
	uc := controller.NewUserController(getDbSession())
	sc := controller.NewSensorController(getDbSession())
	vc := controller.NewVirtSensorController(getDbSession())

	router := httprouter.New()

	//Project Management Routings
	router.Handle("GET", "/projects", pc.Project)
	router.Handle("POST", "/addproject", pc.AddProject)
	router.Handle("GET", "/projects/:projectName/selectsensors", pc.SelectSensors)
	router.Handle("POST", "/projects/:projectName/addvirtualsensors", pc.AddSensorToProject)
	router.Handle("GET", "/projects/:projectName/dashboard", pc.Dashboard)
	router.Handle("GET", "/projects/:projectName/terminate", pc.TerminateProject)
	/*router.Handle("GET", "/sendvirtualsensordetails/", handle)
	router.Handle("POST", "/projects/:projectName/addvirtualsensors", handle)
	router.Handle("GET", "/projects/:projectName/startvirtsensor/:vsensorid", handle)
	router.Handle("GET", "/projects/:projectName/stopvirtsensor/:vsensorid", handle)
	router.Handle("GET", "/projects/:projectName/terminatevirtsensor/:vsensorid", handle)*/

	// User Management Routings
	router.Handle("GET", "/", uc.Index)
	router.Handle("GET", "/logout", uc.Logout)
	router.Handle("POST", "/users/login", uc.Login)
	router.Handle("POST", "/users/signup", uc.Signup)
	router.Handle("PUT", "/users/updateuser/:username", uc.UpdateUser)
	router.Handle("DELETE", "/users/deleteuser/:username", uc.DeleteUser)
	router.Handle("GET", "/time", uc.Datahandler)
	router.Handle("GET", "/gettime", uc.HandlerGetData)

	//Physical Sensor Management Routings
	router.Handle("POST", "/sensors", sc.AddSensor)
	router.Handle("PUT", "/sensors/updatesensor/:sensorName", sc.UpdateSensor)
	router.Handle("DELETE", "/sensors/deletesensor/:sensorName", sc.DeleteSensor)
	router.Handle("GET", "/sensors/getsensor/:sensorName", sc.GetSensorByName)
	router.Handle("POST", "/authenticatesenor", sc.AuthenticateSensor)
	router.Handle("GET", "/getsensorbyowner/:ownerId", sc.GetSensorByOwner)
	router.Handle("GET", "/getallsharedsensors", sc.GetAllSharedSensors)
	//router.Handle("POST", "/getsensordata",sc.GetSensorData)
	//router.Handle("PUT", "/sensors/togglesensor/:sensorName", ToggelSensor)
	/*router.Handle("GET", "/getsensorbyowner/:ownerId", GetSensorById)*/

	// Virtual Sensor Management Routings
	router.Handle("POST", "/virtsensors", vc.AddVirtSensor)
	router.Handle("PUT", "/removevirtsensors/:vSensorName", vc.RemoveVirtSensor)
	router.Handle("GET", "/stopvirtsensors/:vSensorName", vc.StopVirtSensor)
	router.Handle("GET", "/resumevirtsensors/:vSensorName", vc.ResumeVirtSensor)
	router.Handle(("GET"), "/projects/:projectName/sensors/:vSensorName", vc.GetVSensorDetails)

	//fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	router.NotFound = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":8081", router))
}
