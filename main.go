package main

import (
	"github.com/leroxyl/greenbone/database"
	"github.com/leroxyl/greenbone/web"
)

func main() {
	// initialize database
	cm := database.NewComputerManager()

	// start server
	web.Run(cm)
}
