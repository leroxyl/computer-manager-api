package cmd

import (
	"github.com/leroxyl/computer-manager-api/internal/adapter/client"
	"github.com/leroxyl/computer-manager-api/internal/adapter/storage"
	"github.com/leroxyl/computer-manager-api/internal/adapter/web"
	"github.com/leroxyl/computer-manager-api/internal/config"

	log "github.com/sirupsen/logrus"
)

func Run() {
	log.Infof("Starting Computer Manager API")

	// read configuration
	conf := config.Config{}
	conf.Load()

	// initialize notification service client
	nsc := client.NewNotificationServiceClient(conf.ClientConfig)

	// initialize database
	dm := storage.NewDatabaseManager(conf.DatabaseConfig, nsc.NotifyAdmin)

	// set up server
	r := web.NewServer(dm)

	// start server (blocks)
	r.Run()
}
