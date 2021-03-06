package main

import (
	"flag"
	"forum/internal/api"
	"log"
	"strconv"
)

var (
	port int
	db   string
	//createNewDb bool
)

func init() {
	flag.IntVar(&port, "port", 8081, "Specify the port to listen to.")
	flag.StringVar(&db, "db", "./dataBase.db", "Specify path to database")
	//flag.BoolVar(&createNewDb, "createDB", false, "Specify whether to create a new database")
}

// @title forum API
// @version 1.0
// @description This is a sample service for managing orders
// @termsOfService http://swagger.io/terms/
// @host localhost:8081
// @BasePath /api

func main() {
	flag.Parse()
	log.Println("It works") //добавить инфо логгер

	//server instance initialization
	config := api.NewConfig(strconv.Itoa(port), db)
	server := api.New(config)

	//start server
	log.Fatal(server.Start())
}
