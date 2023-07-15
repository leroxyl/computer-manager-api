package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/leroxyl/greenbone/database"
)

type ComputerManager interface {
	Create(*database.Computer) error
	Read(*database.Computer) error
}

// Run initializes all endpoints and starts the server.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func Run(cm ComputerManager) {
	r := gin.Default()

	r.POST("/computers", createComputer(cm))
	r.GET("/computers/:mac", readComputer(cm))

	// listen and serve
	err := r.Run()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func createComputer(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		computer := &database.Computer{}
		err := c.BindJSON(computer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = cm.Create(computer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}

func readComputer(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		computer := &database.Computer{
			MACAddr: c.Param("mac"),
		}

		err := cm.Read(computer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}
