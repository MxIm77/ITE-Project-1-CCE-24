package main

import (
	PhoeniciaDigitalServer "Phoenicia-Digital-Base-API/base/server"
)

func main() {

	//	if Postgres Database Not In use comment out
	// defer PhoeniciaDigitalDatabase.Postgres.DB.Close()

	// if MongoDB Database Not In Use comment out
	// defer PhoeniciaDigitalDatabase.Mongo.Client.Disconnect(context.Background())

	PhoeniciaDigitalServer.StartServer()
}
