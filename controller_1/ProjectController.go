package controller

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"model_1"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type (
	ProjectController struct {
		db *sql.DB
	}
)

var project = template.Must(template.ParseFiles(
	"templates/projects.tmpl",
	"templates/physicalsensorlist.tmpl",
	"templates/projectlist.tmpl",
))

var psensor = template.Must(template.ParseFiles(
	"templates/addVirtualSensor.tmpl",
	"templates/sensorlist.tmpl",
))

var projectdashboard = template.Must(template.ParseFiles(
	"templates/dashboard.tmpl",
	"templates/vtable.tmpl",
))

func NewProjectController(db *sql.DB) *ProjectController {
	return &ProjectController{db}
}

func (pc ProjectController) Project(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	type ProjectList struct {
		ProjectName string
		ProjectLink string
	}

	type SensorList struct {
		SensorName string
		SensorLink string
	}

	type Page struct {
		Plist *[]ProjectList
		Slist *[]SensorList
	}

	page := Page{}
	pcount := 0
	scount := 0
	cookie, _ = r.Cookie("logged-in")

	// Display the Projects for a User
	rows, err := pc.db.Query("Select count(*) from project_metadata where project_owner=? and project_status!=?", cookie.Value, "Terminated")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		rows.Scan(&pcount)
	}
	fmt.Println("Number of Projects:", pcount)
	projectList := make([]ProjectList, pcount)

	// Get Projects
	rows, err = pc.db.Query("Select project_name from project_metadata where project_owner=? and project_status!=?", cookie.Value, "Terminated")
	if err != nil {
		fmt.Println(err)
	}

	index := 0
	for rows.Next() {
		rows.Scan(&projectList[index].ProjectName)
		projectList[index].ProjectLink = fmt.Sprintf("/projects/%s/dashboard", projectList[index].ProjectName)
		fmt.Println(err)
		index++
	}

	fmt.Println(projectList)
	page.Plist = &projectList

	// Display Sensors for a User
	rows, err = pc.db.Query("Select count(*) from physical_sensor_metadata where owner_id=?", cookie.Value)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		rows.Scan(&scount)
	}

	sensorList := make([]SensorList, scount)

	// Get Sensors
	rows, err = pc.db.Query("Select sensor_name from physical_sensor_metadata where owner_id=?", cookie.Value)
	if err != nil {
		fmt.Println(err)
	}

	index = 0
	for rows.Next() {
		rows.Scan(&sensorList[index].SensorName)
		sensorList[index].SensorLink = fmt.Sprintf("/sensors/getsensor/%s", sensorList[index].SensorLink)
		fmt.Println(err)
		index++
	}
	page.Slist = &sensorList
	if err := project.Execute(w, page); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (pc ProjectController) AddProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	cookie, _ = r.Cookie("logged-in")
	project := model.Project{}

	//Random Generator for the User Id generation
	randId, _ := exec.Command("uuidgen").Output()
	project.ProjectId = string(randId)
	project.ProjectName = r.FormValue("projectName")
	project.ProjectDesc = r.FormValue("projectDesc")
	project.ProjectOwner = cookie.Value
	project.ProjectStatus = "New"
	fmt.Println(project)
	//Insert User in DB
	stmt, err := pc.db.Prepare("Insert project_metadata SET project_id=?,project_name=?,project_desc=?,project_owner=?,project_status=?")
	res, err := stmt.Exec(project.ProjectId, project.ProjectName, project.ProjectDesc, project.ProjectOwner, project.ProjectStatus)
	if err != nil {
		fmt.Errorf("Error in registering the Project: %v", err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Print(err)
	}
	redirectUrl := fmt.Sprintf("/projects/%s/selectsensors", project.ProjectName)
	if affect == 1 {
		projectcookie := &http.Cookie{
			Name:     "project-id",
			Value:    project.ProjectId,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, projectcookie)
		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
	}

}

func (pc ProjectController) SelectSensors(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	type Physical struct {
		SensorName string
		SensorType string
		Latitude   string
		Longitude  string
	}

	type Prj struct {
		ProjectURL string
	}

	type Page struct {
		PhysicalList *[]Physical
		Project      *Prj
	}

	page := Page{}
	prj := Prj{}

	//prj.ProjectURL = fmt.Sprintf("/projects/%s/addvirtualsensors", ps.ByName("projectName"))
	prj.ProjectURL = ps.ByName("projectName")
	fmt.Println(prj.ProjectURL)
	page.Project = &prj

	pcount := 0
	// Physical Sensor Display
	rows, err := pc.db.Query("Select count(*) from physical_sensor_metadata where is_shared=true and is_activated=true")
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	for rows.Next() {
		rows.Scan(&pcount)
	}
	fmt.Println(pcount)
	rows, err = pc.db.Query("Select sensor_name, sensor_type, latitude, longitude from physical_sensor_metadata where is_shared=true and is_activated=true")
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	phySensorList := make([]Physical, pcount)
	index := 0
	for rows.Next() {
		rows.Scan(&phySensorList[index].SensorName, &phySensorList[index].SensorType, &phySensorList[index].Latitude, &phySensorList[index].Longitude)
		fmt.Println(err)
		fmt.Println(phySensorList[index])
		index++
	}
	fmt.Println(index)
	page.PhysicalList = &phySensorList
	err = psensor.Execute(w, page)
	if err != nil {
		fmt.Println(err)
	}
}

func (pc ProjectController) AddSensorToProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	type SensorList struct {
		Sensor  []string
		VSensor []string
	}

	type Configuration struct {
		NumberOfSensors int      `json:"numsensor"`
		SensorId        []string `json:"sensorid"`
		Status          []bool   `json:"status"`
		PhySensorName   []string `json:"physensorname"`
		Frequency       []int    `json:"freq"`
	}

	cookie, _ := r.Cookie("logged-in")
	prjcookie, _ := r.Cookie("project-id")
	projectName := ps.ByName("projectName")
	fmt.Println(projectName)
	var SensorString string
	var VSensorString string

	sensorList := SensorList{}
	vsensor := model.VirtSensor{}

	fmt.Println(r.URL)
	r.ParseForm()
	//fmt.Println(r.Form)
	fmt.Println(r.Form["sensor[]"])
	fmt.Println(len(r.Form["sensor[]"]))

	for k, v := range r.Form {
		if k == "sensor[]" {
			for _, val := range v {
				SensorString += val + ","
			}
		} else {
			for _, val := range v {
				if val != "" {

					VSensorString += val + ","
				}
			}
		}
	}

	fmt.Println("Sensor:", SensorString)
	fmt.Println("Vsensor:", VSensorString)
	sensorList.Sensor = strings.Split(SensorString, ",")
	sensorList.Sensor = sensorList.Sensor[0 : len(sensorList.Sensor)-1]
	sensorList.VSensor = strings.Split(VSensorString, ",")
	sensorList.VSensor = sensorList.VSensor[0 : len(sensorList.VSensor)-1]
	fmt.Println("Sensor:", sensorList.Sensor)
	fmt.Println("Vsensor:", sensorList.VSensor)

	vsensorid := make([]string, len(sensorList.Sensor))
	vsensorStatus := make([]bool, len(sensorList.Sensor))
	vsensorFreq := make([]int, len(sensorList.Sensor))
	//Insert into virtualsensor-metadata
	index := 0
	for index < len(sensorList.Sensor) {

		randId, _ := exec.Command("uuidgen").Output()
		vsensor.VirtSensorId = string(randId)
		vsensorid[index] = vsensor.VirtSensorId
		vsensor.VirtSensorName = sensorList.VSensor[index]
		vsensor.VirtSensorUserId = cookie.Value
		vsensor.ProjectId = prjcookie.Value
		vsensor.PhysicalSensorName = sensorList.Sensor[index]
		vsensor.IsEnabled = false
		vsensor.DataFreqinSec = 5
		vsensorFreq[index] = vsensor.DataFreqinSec
		vsensor.RegisterTime = time.Now().UTC()
		vsensor.Status = 0
		vsensorStatus[index] = false

		stmt, err := pc.db.Prepare("INSERT INTO `virtual_sensor_metadata` VALUES(?,?,?,?,?,?,?,?,?,?)")
		_, err = stmt.Exec(vsensor.VirtSensorId, vsensor.VirtSensorUserId, vsensor.PhysicalSensorName, vsensor.ProjectId, vsensor.VirtSensorName, vsensor.IsEnabled, vsensor.DataFreqinSec, 0, vsensor.RegisterTime, vsensor.Status)
		if err != nil {
			fmt.Println(err)
		}

		index++
	}

	stmt, err := pc.db.Prepare("INSERT INTO `docker_metadata` (project_id) VALUES(?)")
	_, err = stmt.Exec(vsensor.ProjectId)
	if err != nil {
		fmt.Println(err)
	}

	config := Configuration{}
	config.NumberOfSensors = len(sensorList.Sensor)
	config.Frequency = vsensorFreq
	config.PhySensorName = sensorList.Sensor
	config.SensorId = vsensorid
	config.Status = vsensorStatus
	jsonConfig, _ := json.Marshal(config)

	go func() {
		flag := true
		transport := http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Duration(1*time.Second))
			},
		}
		client := http.Client{
			Transport: &transport,
		}

		var port string
		err := pc.db.QueryRow("SELECT docker_id FROM docker_metadata WHERE project_id=?", vsensor.ProjectId).Scan(&port)
		if err != nil {
			fmt.Println(err)
		}
		strport := port + ":9999"
		go sh.Command("docker", "run", "--publish", strport, "--name", projectName, "--rm", "virtsensor", "&").Run()

		for flag == true {

			resp, err := client.Get(fmt.Sprintf("http://localhost:%s/isalive", port))
			if err != nil {
				fmt.Println("Sleeping...")
				time.Sleep(time.Duration(1000 * time.Millisecond))
			} else {
				if resp.StatusCode == 200 {
					flag = false
					fmt.Println(config)
					buf := new(bytes.Buffer)
					err = binary.Write(buf, binary.BigEndian, &jsonConfig)
					configurl := fmt.Sprintf("http://localhost:%s/getconfiguration", port)
					req, _ := http.NewRequest("POST", configurl, buf)
					req.Header.Set("Content-Type", "application/json")
					_, err = client.Do(req)
				}
			}
		}
	}()

	redirectUrl := fmt.Sprintf("/projects/%s/dashboard", projectName)
	fmt.Println(redirectUrl)
	http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
}

/*func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(1*time.Second))
}*/

func (pc ProjectController) Dashboard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	type Virtual struct {
		Usage          int
		Status         bool
		PhysicalSensor string
		SensorName     string
		ProjectName    string
		Is_Enabled     bool
		DataFrequency  int
		Created        time.Time
	}

	vcount := 0

	cookie, _ := r.Cookie("logged-in")
	projectName := ps.ByName("projectName")

	var projectId string

	rows, err := pc.db.Query("Select project_id from project_metadata where project_name=? and project_owner=?", projectName, cookie.Value)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	for rows.Next() {
		rows.Scan(&projectId)
	}

	projectcookie := &http.Cookie{
		Name:     "project-id",
		Value:    projectId,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, projectcookie)
	fmt.Println("Project ID:", strings.Trim(projectcookie.Value, "\n"))
	fmt.Println("Cookie Value:", cookie.Value)

	// Virtual Sensor Display for User
	rows, err = pc.db.Query("Select count(*) from virtual_sensor_metadata where enduser_id=? and project_id=?", cookie.Value, strings.Trim(projectcookie.Value, "\n"))
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	for rows.Next() {
		rows.Scan(&vcount)
	}
	fmt.Println("Virtual Sensor Count:", vcount)

	rows, err = pc.db.Query("Select psensor_name, usage_counter_in_mins, vsensor_name, status, is_enabled , data_frequency_in_sec, created from virtual_sensor_metadata where enduser_id=? and project_id=?", cookie.Value, strings.Trim(projectcookie.Value, "\n"))
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	virt := make([]Virtual, vcount)
	index := 0
	for rows.Next() {
		rows.Scan(&virt[index].PhysicalSensor, &virt[index].Usage, &virt[index].SensorName, &virt[index].Status, &virt[index].Is_Enabled, &virt[index].DataFrequency, &virt[index].Created)
		virt[index].ProjectName = projectName
		fmt.Println(err)
		index++
	}
	err = projectdashboard.Execute(w, virt)
	if err != nil {
		fmt.Println(err)
	}
}

func (pc ProjectController) TerminateProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	cookie, _ := r.Cookie("logged-in")
	projectName := ps.ByName("projectName")
	stmt, err := pc.db.Prepare("update project_metadata set project_status='Terminated' where project_name=? and project_owner=?")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(projectName, cookie.Value)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/projects", http.StatusMovedPermanently)
}
