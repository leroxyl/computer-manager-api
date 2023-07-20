package cmd

import (
	"github.com/leroxyl/computer-manager-api/internal/adapter/storage"
	"github.com/leroxyl/computer-manager-api/internal/adapter/web"

	log "github.com/sirupsen/logrus"
)

func Run() {
	log.Infof("Starting Computer Manager API")

	// initialize database
	cm := storage.NewDatabaseManager()

	// set up server
	r := web.NewServer(cm)

	// start server (blocks)
	r.Run()
}
