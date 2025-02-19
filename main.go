package main

import (
	"app/config"
	"app/database"
	"app/server"
	"flag"
)

func setup() {
	flag.StringVar(&config.Config.Host, "host", "0.0.0.0", "host for the redis server")
	flag.UintVar(&config.Config.Port, "port", 7379, "port for the redis server")
	flag.UintVar(&config.Config.MaxClients, "max-clients", 30000, "max number of simultaneously connected clients")
	flag.Parse()
}

func main() {
	setup()

	db := database.NewDatabase()
	s, _ := server.NewServer(&config.Config, db)
	s.RunServer()
}
