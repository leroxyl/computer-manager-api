package cmd

import (
	"github.com/leroxyl/computer-manager-api/internal/adapter/storage"
	"github.com/leroxyl/computer-manager-api/internal/adapter/web"
)

func Run() {
	// initialize database
	cm := storage.NewComputerManager()

	// set up server
	r := web.NewServer(cm)

	// start server (blocks)
	r.Run()
}
