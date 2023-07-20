package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leroxyl/computer-manager-api/internal/entity"
	"github.com/leroxyl/computer-manager-api/internal/usecase"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	router          *gin.Engine
	computerManager usecase.ComputerManager
}

// NewServer initializes a new server instance and sets up all endpoints
func NewServer(cm usecase.ComputerManager) *Server {
	r := gin.Default()

	server := &Server{
		router:          r,
		computerManager: cm,
	}

	// for this application we don't need any proxy headers
	_ = r.SetTrustedProxies(nil)

	r.POST("/computers", server.createComputer())
	r.GET("/computers/:mac", server.readComputer())
	r.PUT("/computers/:mac", server.updateComputer())
	r.DELETE("/computers/:mac", server.deleteComputer())
	r.GET("/computers", server.readAllComputers())
	r.GET("/employees/:abbr/computers", server.readAllComputersForEmployee())

	return server
}

// Run starts the server.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (s *Server) Run() {
	err := s.router.Run()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (s *Server) createComputer() gin.HandlerFunc {
	return func(c *gin.Context) {
		computer := entity.Computer{}
		err := c.BindJSON(&computer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = s.computerManager.Create(computer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}

func (s *Server) readComputer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mac := getMACAddress(c)

		computer, err := s.computerManager.Read(mac)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}

func (s *Server) updateComputer() gin.HandlerFunc {
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

		err = s.computerManager.Update(computer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, computer)
	}
}

func (s *Server) deleteComputer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mac := getMACAddress(c)

		err := s.computerManager.Delete(mac)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.String(http.StatusOK, "entry deleted")
	}
}

func (s *Server) readAllComputers() gin.HandlerFunc {
	return func(c *gin.Context) {
		computers, err := s.computerManager.ReadAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, computers)
	}
}

func (s *Server) readAllComputersForEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		abbr := c.Param("abbr")

		computers, err := s.computerManager.ReadAllForEmployee(abbr)
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
