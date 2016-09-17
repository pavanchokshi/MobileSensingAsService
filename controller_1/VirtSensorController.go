package controller

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/julienschmidt/httprouter"
	"log"
	//"model"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type (
	VirtSensorController struct {
		db *sql.DB
	}
)

func NewVirtSensorController(db *sql.DB) *VirtSensorController {
	return &VirtSensorController{db}
}

var ExposedPort = 6060

var vsensordetails = template.Must(template.ParseFiles(
	"templates/vsensordetails.tmpl",
	"templates/vsensordetailsform.tmpl",
))

func ComputeResourceAvailability() {

	// Lazy Approach !!

}

// Multiple States of Control for the Physical Sensor
// Physical Sensor Owner has not shared the sensor: Returns -111
// Physical Sensor Owner has shared the sensor but no Physical Controller Exists: Returns 0
// Physical Sensor Owner has shared the sensor also Physical Controller Exists: Returns 1
// Physical Sensor Owner has shared the sensor but the current status of Sensor is Disabled : Returns 400
func (vc VirtSensorController) IsPhysicalControllerPresent(sensorName string) int {

	var sensor_name string
	var is_shared bool
	var owner_enabled bool
	var countVirtualSensors int
	// Queries the Sensor Metadata to fetch the sensor object
	err := vc.db.QueryRow("Select sensor_name, is_shared, owner_enabled from physical_sensor_metadata where sensor_name=?", sensorName).Scan(&sensor_name, &is_shared, &owner_enabled)
	if err != nil {
		log.Fatal(err)
	}

	if is_shared {

		if owner_enabled {

			err = vc.db.QueryRow("Select count_v_sensor from sensor_stat where sensor_name=?", sensorName).Scan(&countVirtualSensors)
			if err != nil {
				log.Fatal(err)
			}
			if countVirtualSensors > 0 {
				return 1
			} else {
				return 0
			}
		} else {

			return 400
		}

	} else {

		return -1
	}
}

// Status of Sensor:
// 1-Active
// 2-Stopped
// 3-Terminated

func (vc VirtSensorController) AddVirtSensor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//vsensor := model.VirtSensor{}
	cookie, _ = r.Cookie("logged-in")
	//temp := vc.IsPhysicalControllerPresent("Dhruv")

	//Random Generator for the User Id generation
	randId, _ := exec.Command("uuidgen").Output()

	/*vsensor.VirtSensorId = string(randId)
	vsensor.VirtSensorName = r.FormValue("virtSensorName")
	vsensor.VirtSensorUserId = cookie.Value
	vsensor.PhysicalSensorName = r.FormValue("phsense")
	vsensor.DataFreqinSec = 1
	vsensor.RegisterTime = time.Now().UTC()
	vsensor.IsEnabled = false
	vsensor.UsageCounterinMins = 0
	vsensor.Status = 0
	fmt.Println("Sensor Name:", vsensor.PhysicalSensorName)
	fmt.Println("Adding Sensor:", vsensor)*/
	//Insert User in DB
	stmt, err := vc.db.Prepare("Insert INTO virtual_sensor_metadata VALUES (?,?,?,?,?,?,?,?,?)")
	fmt.Println(stmt)
	if err != nil {
		fmt.Println(err)
		_, err = stmt.Exec(string(randId), cookie.Value, r.FormValue("phsense"), r.FormValue("virtSensorName"), 0, 1, 0, time.Now().UTC(), 0)
		fmt.Println(err)
	}
	http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
}

func (vc VirtSensorController) RemoveVirtSensor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	vSensorName := ps.ByName("vSensorName")
	var vSensorId string
	var status int

	// Fetch sensor id from the VSensor Metadata table
	row, err := vc.db.Query("Select vsensor_id from virtual_sensor_metadata where vsensor_name=?", vSensorName)
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		row.Scan(&vSensorId)
		if err != nil {
			log.Fatal(err)
		}

	}

	status = 3 // Terminated
	stmt, err := vc.db.Prepare("update virtual_sensor_metadata set status=? where vsensor_id=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(status, vSensorId)
	if err != nil {
		log.Fatal(err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affect)

}

func (vc VirtSensorController) StopVirtSensor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	vSensorName := ps.ByName("vSensorName")
	var vSensorId string
	var status int

	// Fetch sensor id from the VSensor Metadata table
	row, err := vc.db.Query("Select vsensor_id from virtual_sensor_metadata where vsensor_name=?", vSensorName)
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		row.Scan(&vSensorId)
		if err != nil {
			log.Fatal(err)
		}

	}

	status = 2 // Stopped
	stmt, err := vc.db.Prepare("update virtual_sensor_metadata set status=? where vsensor_id=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(status, vSensorId)
	if err != nil {
		log.Fatal(err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affect)

}

func (vc VirtSensorController) ResumeVirtSensor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	vSensorName := ps.ByName("vSensorName")
	var vSensorId string
	var status int
	ExposedPort += 1

	// Fetch sensor id from the VSensor Metadata table
	row, err := vc.db.Query("Select vsensor_id from virtual_sensor_metadata where vsensor_name=?", vSensorName)
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		row.Scan(&vSensorId)
		if err != nil {
			log.Fatal(err)
		}

	}

	status = 1 // Started
	stmt, err := vc.db.Prepare("update virtual_sensor_metadata set status=? where vsensor_id=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(status, vSensorId)
	if err != nil {
		log.Fatal(err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affect)
	port := strconv.Itoa(ExposedPort)
	strport := fmt.Sprintf("%s:9000", port)
	go sh.Command("docker", "run", "--publish", strport, "--name", vSensorName, "--rm", "virtsensor", "&").Run()
}

func (vc VirtSensorController) GetVSensorDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	prjcookie, _ := r.Cookie("project-id")
	vSensorName := ps.ByName("vSensorName")

	type Details struct {
		SensorId    string
		SensorName  string
		PSensorName string
		URL         string
		Freq        int
	}
	details := Details{}

	err := vc.db.QueryRow("SELECT vsensor_id, psensor_name,data_frequency_in_sec FROM virtual_sensor_metadata WHERE project_id=? and vsensor_name=?", prjcookie.Value, vSensorName).Scan(&details.SensorId, &details.PSensorName, &details.Freq)
	if err != nil {
		fmt.Println(err)
	}

	var port string
	err = vc.db.QueryRow("SELECT docker_id FROM docker_metadata WHERE project_id=?", prjcookie.Value).Scan(&port)
	if err != nil {
		fmt.Println(err)
	}

	details.SensorName = vSensorName
	details.URL = "http://localhost:" + port
	if err := vsensordetails.Execute(w, details); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
