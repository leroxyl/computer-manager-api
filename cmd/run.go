package cmd

import (
	"github.com/leroxyl/greenbone/internal/adapter/storage"
	"github.com/leroxyl/greenbone/internal/adapter/web"
)

func Run() {
	// initialize database
	cm := storage.NewComputerManager()

	// start server
	web.Run(cm)
}
