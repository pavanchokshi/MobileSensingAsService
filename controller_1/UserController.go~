package controller

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"math/rand"
	"model_1"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type (
	UserController struct {
		db *sql.DB
	}
)

var cookie *http.Cookie

var index = template.Must(template.ParseFiles(
	"templates/login.html",
))

var page = `<!DOCTYPE html>
<html lang="en">

<head>
<title>Charts - Bootstrap Admin Template</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <link href="./static/css/bootstrap.min.css" rel="stylesheet">
    <link href="./static/css/bootstrap-responsive.min.css" rel="stylesheet">
    <link href="http://fonts.googleapis.com/css?family=Open+Sans:400italic,600italic,400,600" rel="stylesheet">
    <link href="./static/css/font-awesome.css" rel="stylesheet">
    <link href="../static/css/style.css" rel="stylesheet">
    <script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.3.2/jquery.min.js"></script>
  <script type="text/javascript">
  window.onload = function () {

    var dps = []; // dataPoints
    
    var chart = new CanvasJS.Chart("chartContainer",{
      title :{
        text: "Live Random Data"
      },      
      axisX:{
          title: "Time",
          valueFormatString:  "hh:mm:ss" ,
          interval: 1,
          intervalType: "second"
        
      },
      data: [{
        type: "line",
        dataPoints: dps 
      }],
    });

    var xVal = new Date();
    var yVal = 0; 
    var updateInterval = 100;
    
    var dataLength = 7; // number of dataPoints visible at any point


    var updateChart = function () {
      
      $.post("http://localhost:8080/gettime", "", function(data, status) {
                  $("#output").empty();
                  $("#output").append(data);
                  xVal = new Date()
                  yVal = parseInt(data);
                  alert(data);
            });

        dps.push({
          x: xVal,
          y: yVal
        });
      
        xVal++;
      if (dps.length > dataLength)
      {
        dps.shift();        
      }
      
      chart.render();   

    };

    // generates first set of dataPoints
    updateChart(dataLength); 
    // update chart after specified time. 
    var myVar = setInterval(function(){updateChart()}, 1000); 
    $(document).ready(function () {
                $("#sensor_button").click(function(){
                        clearInterval(myVar);
                        if(document.getElementById("sensor_button").getAttribute("class") == "btn btn-danger"){
                            document.getElementById("sensor_button").setAttribute("class", "btn btn-success");     
                            document.getElementById("sensor_button").innerHTML = "Start"
                        }else{
                            document.getElementById("sensor_button").setAttribute("class", "btn btn-danger");     
                            document.getElementById("sensor_button").innerHTML = "Stop"
                            myVar = setInterval(function(){updateChart()}, 1000); 
                        }
                });
                        $("#output").append("Waiting for system time..");
                      //myVar = setInterval(function(){updateChart()}, 1000);
                   });
                }
            function myStopFunction() {
            clearInterval(myVar);
            if(document.getElementById("sensor_button").getAttribute("class") == "btn btn-danger"){
                document.getElementById("sensor_button").setAttribute("class", "btn btn-success");     
                document.getElementById("sensor_button").innerHTML = "Start"
            }else{
                document.getElementById("sensor_button").setAttribute("class", "btn btn-danger");     
                document.getElementById("sensor_button").innerHTML = "Stop"
                myVar = setInterval(function(){updateChart()}, 10000); 
            }
            
        }
  </script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/canvasjs/1.7.0/jquery.canvasjs.min.js"></script>
</head>
<body>
    <div class="navbar navbar-fixed-top">
        <div class="navbar-inner">
            <div class="container">
                <a class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse"><span
                    class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span>
                </a><a class="brand" href="index.html">Mobile Sensor Cloud </a>
                <div class="nav-collapse">
                    <ul class="nav pull-right">
                        <li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown"><i
                            class="icon-user"></i> Hi, Pavan <b class="caret"></b></a>
                            <ul class="dropdown-menu">
                                <li><a href="javascript:;">Profile</a></li>
                                <li><a href="javascript:;">Logout</a></li>
                            </ul>
                        </li>
                    </ul>
                </div>
                <!--/.nav-collapse -->
            </div>
            <!-- /container -->
        </div>
        <!-- /navbar-inner -->
    </div>
    <!-- /navbar -->
    <div class="subnavbar">
        <div class="subnavbar-inner">
            <div class="container">
                <ul class="mainnav">
                    <li><a href="index.html"><i class="icon-dashboard"></i><span>Dashboard</span> </a>
                    </li>
                    <li><a href="#myModal" role="button" data-toggle="modal"><i class="icon-dashboard"></i><span>App Tour</span>
                    </a></li>
                    <li class="active"><a href="charts.html"><i class="icon-bar-chart"></i><span>Charts</span> </a>
                    </li>
                </ul>
            </div>
            <!-- /container -->
        </div>
        <!-- /subnavbar-inner -->
    </div>
    <!-- /subnavbar -->
    <div class="main">
        <div class="main-inner">
            <div class="container">
                <div class="row">
                    <div class="span6">
                       <!-- /widget -->
                        <div class="widget">
                            <div class="widget-header">
                                <i class="icon-bar-chart"></i>
                                <h3>
                                    Live Graph</h3>
                            </div>
                            <!-- /widget-header -->
                            <div class="widget-content">
                                <div id="chartContainer" style="height: 300px; width:100%;">
                            </div>
                                <!-- /pie-chart -->
                            </div>
                            <!-- /widget-content -->
                        </div>
                        <!-- /widget -->
                    </div>
                    <!-- /span6 -->
                    <div class="span6">
                        <div class="widget">
                            <div class="widget-header">
                                <i class="icon-bar-chart"></i>
                                <h3>
                                    Control</h3>
                            </div>
                            <!-- /widget-header -->
                            <div class="widget-content">
                                <button class="btn btn-danger" id="sensor_button" onclick="myStopFunction()">Stop</button>
                            </div>
                                <!-- /pie-chart -->
                            </div>
                            <!-- /widget-content -->
                        </div>
                    </div>
                    <!-- /span6 -->
                </div>
                <!-- /row -->
            </div>
            <!-- /container -->
        </div>
        <!-- /main-inner -->
    </div>
    <!-- /main -->
</body>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/canvasjs/1.7.0/jquery.canvasjs.min.js"></script>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/flot/0.8.3/excanvas.js"></script>
    <script src="./static/js/chart.min.js" type="text/javascript"></script>
    <script src="./static/js/bootstrap.js"></script>
    <script src="./static/js/base.js"></script>
</html>`

func NewUserController(db *sql.DB) *UserController {
	return &UserController{db}
}

func (uc UserController) Datahandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, page)
}

func (uc UserController) HandlerGetData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Fprint(w, time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST"))
	myrand := random(1, 15)
	fmt.Fprintf(w, strconv.Itoa(myrand))
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (uc UserController) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	cookie, err := r.Cookie("logged-in")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
			//HttpOnly: true,
		}
	}
	fmt.Println(cookie)
	if err := index.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc UserController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { // User Login

	var passDb string
	var roleDb string
	var userId string

	user := model.FormUserLogin{}

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	//Creates and stores the password as password hash
	//Uses SHA-256 as the hashing algorithm
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum(nil))

	// Queries the user_db database to fetch the role and password of the username
	row, err := uc.db.Query("Select user_db_id,password,role from user_db where username=?", user.Username)
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		row.Scan(&userId, &passDb, &roleDb)
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println(userId)
	if user.Password == passDb {
		cookie = &http.Cookie{
			Name:     "logged-in",
			Value:    userId,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		//http.RedirectHandler("/dashboard", http.StatusMovedPermanently)
		http.Redirect(w, r, "/projects", http.StatusMovedPermanently)
		//fmt.Fprint(w, "Login Successful")
	} else {
		fmt.Fprint(w, "Bad Login Request")
	}

}

func (uc UserController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	cookie = &http.Cookie{
		Name:   "logged-in",
		Value:  "0",
		MaxAge: -1,
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	//fmt.Println(cookie)
}

func (uc UserController) Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { // User Signup

	user := model.User{}
	fmt.Println(r)
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.EmailId = r.FormValue("email")
	user.Role = r.FormValue("user_type")

	/*// JSON to Object Marshalling
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&user)
	if err != nil {
		fmt.Errorf("Error in decoding the JSON: %v", err)
	}*/

	//Creates and stores the password as password hash
	//Uses SHA-1 as the hashing algorithm
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum(nil))

	//Random Generator for the User Id generation
	randId, _ := exec.Command("uuidgen").Output()
	user.UserId = string(randId)

	//Insert User in DB
	stmt, err := uc.db.Prepare("Insert user_db SET user_db_id=?,username=?,emailid=?,password=?,role=?")
	res, err := stmt.Exec(user.UserId, user.Username, user.EmailId, user.Password, user.Role)
	if err != nil {
		fmt.Errorf("Error in registering the user: %v", err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if affect == 1 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	username := ps.ByName("username")
	user := model.User{}
	var userId string

	// JSON to Object Marshalling
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&user)
	if err != nil {
		fmt.Errorf("Error in decoding the JSON: %v", err)
	}

	//Creates and stores the password as password hash
	//Uses SHA-1 as the hashing algorithm
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum(nil))

	row, err := uc.db.Query("Select user_db_id from user_db where username=?", username)
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		row.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}

	}

	stmt, err := uc.db.Prepare("update user_db set emailid=?,password=? where user_db_id=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(user.EmailId, user.Password, userId)
	if err != nil {
		log.Fatal(err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affect)

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	username := ps.ByName("username")
	var userId string
	row, err := uc.db.Query("Select user_db_id from user_db where username=?", username)
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		row.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}

	}

	stmt, err := uc.db.Prepare("delete from user_db where user_db_id=?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(userId)
	if err != nil {
		log.Fatal(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(affect)
}
