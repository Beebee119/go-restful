package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"

	metrodb "github.com/Beebee119/go-restful/metro-rail/db"
	_ "github.com/mattn/go-sqlite3"
)

// var DB *sql.DB

type Train struct {
	ID              int
	DriverName      string
	OperatingStatus bool
	DB              *sql.DB
}
type Station struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
	DB          *sql.DB
}
type Schedule struct {
	ID          int
	TrainID     int
	StaionID    int
	ArrivalTime time.Time
	DB          *sql.DB
}

func (t *Train) Register(container *restful.Container, DB *sql.DB) {
	t.DB = DB
	ws := new(restful.WebService)
	ws.
		Path("/v1/trains").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{train_id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train_id}").To(t.deleteTrain))
	container.Add(ws)
}

func (t *Train) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train_id")
	err := t.DB.QueryRow("select id, driver_name, operating_status from train where id = ?;", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println("Error get Train ID: ", id, err)
		response.AddHeader("Content-Type", "text-plain")
		response.WriteErrorString(http.StatusNotFound, "Train could not be found")
	} else {
		response.WriteEntity(t)
	}
}

func (t *Train) createTrain(request *restful.Request, response *restful.Response) {
	var temp Train
	decoder := json.NewDecoder(request.Request.Body)
	err := decoder.Decode(&temp)
	if err != nil {
		log.Println("Error get data from request body")
	}
	defer request.Request.Body.Close()
	log.Printf("Driver Name:%s, Operating status:%v\n", temp.DriverName, temp.OperatingStatus)
	statement, _ := t.DB.Prepare("INSERT INTO train (driver_name, operating_status) VALUES (?, ?)")
	if err != nil {
		log.Println("Error prepare statement")
	} else {
		log.Println("Successfully prepare statement")
	}
	result, err := statement.Exec(temp.DriverName, temp.OperatingStatus)
	if err != nil {
		log.Println("Failed to insert train")
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
	newID, _ := result.LastInsertId()
	temp.ID = int(newID)
	response.WriteHeaderAndEntity(http.StatusCreated, temp)
}

func (t *Train) deleteTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train_id")
	statement, err := t.DB.Prepare("delete from train where id = ?;")
	if err != nil {
		log.Println("Error in creating table books")
	} else {
		log.Println("Successfully created table books")
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Println("Failed to insert train")
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
	response.WriteHeader(http.StatusOK)
}

func main() {
	// Connect to database
	DB, err := sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("Failed to create DB Driver")
	}
	metrodb.Initialize(DB)
	wsContainter := restful.NewContainer()
	wsContainter.Router(restful.CurlyRouter{})
	// restful.CurlyRouter allows us to use {train_id} path parameter
	t := Train{}
	t.Register(wsContainter, DB)
	log.Println("start listening on localhost:8000")
	server := &http.Server{
		Addr:    ":8000",
		Handler: wsContainter,
	}
	log.Fatal(server.ListenAndServe())
}
