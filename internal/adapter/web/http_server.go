package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leroxyl/greenbone/internal/entity"
)

type ComputerManager interface {
	Create(entity.Computer) error
	Read(mac string) (entity.Computer, error)
	Update(entity.Computer) error
	Delete(mac string) error
	ReadAll() ([]entity.Computer, error)
	ReadAllForEmployee(abbr string) ([]entity.Computer, error)
}

// Run initializes all endpoints and starts the server.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func Run(cm ComputerManager) {
	r := gin.Default()

	r.POST("/computers", createComputer(cm))
	r.GET("/computers/:mac", readComputer(cm))
	r.PUT("/computers/:mac", updateComputer(cm))
	r.DELETE("/computers/:mac", deleteComputer(cm))
	r.GET("/computers", readAllComputers(cm))
	r.GET("/employees/:abbr/computers", readAllComputersForEmployee(cm))

	// listen and serve
	err := r.Run()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func createComputer(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		computer := entity.Computer{}
		err := c.BindJSON(&computer)
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
		mac := getMACAddress(c)

		computer, err := cm.Read(mac)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}

func updateComputer(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		mac := getMACAddress(c)

		computer := entity.Computer{}
		err := c.BindJSON(&computer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// prevent user from updating MAC address
		if computer.MACAddr != "" && computer.MACAddr != mac {
			c.JSON(http.StatusBadRequest, gin.H{"error": "updating MAC address is not supported"})
			return
		}

		// insert MAC address from URL parameter into Computer instance
		computer.MACAddr = mac

		err = cm.Update(computer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}

func deleteComputer(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		mac := getMACAddress(c)

		err := cm.Delete(mac)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.String(http.StatusOK, "entry deleted")
	}
}

func readAllComputers(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		computers, err := cm.ReadAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, computers)
	}
}

func readAllComputersForEmployee(cm ComputerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		abbr := c.Param("abbr")

		computers, err := cm.ReadAllForEmployee(abbr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, computers)
	}
}

// getMACAddress returns the MAC address from the URL parameter
func getMACAddress(c *gin.Context) string {
	// TODO: here we could check if the value of the URL parameter is a valid MAC address
	//  and return an error if it is invalid
	return c.Param("mac")
}
