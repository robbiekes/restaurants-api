package main

import (
	"ex00/internal/db/dbmanager"
	"ex00/internal/httpserver"
)

func main() {
	conn := dbmanager.DBConnect()
	httpserver.RequestHandler()
	defer conn.Close()
}

