package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const notificationServiceURL = "http://localhost:8081/api/notify" // TODO externalize URL

type Alarm struct {
	Level                string `json:"level"`
	EmployeeAbbreviation string `json:"employeeAbbreviation"`
	Message              string `json:"message"`
}

// NotifyAdmin sends an HTTP request to an external notification service in order to inform
// an admin about employees with excessive computer demands
func NotifyAdmin(employeeAbbr string, computerCount int64) {
	alarm := Alarm{
		Level:                "warning",
		EmployeeAbbreviation: employeeAbbr,
		Message:              fmt.Sprintf("employee %s has %d computers", employeeAbbr, computerCount),
	}

	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(alarm)
	if err != nil {
		// this should never happen because we are creating the struct instance within this function,
		// so we can assume that it is always serializable.
		// In production code, we would probably handle the error by logging a meaningful message.
		panic(err)
	}

	resp, err := http.Post(notificationServiceURL, "application/json", buffer)
	if err != nil {
		log.Errorf("admin notification failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Errorf("request to admin notification service failed (%d): %q", resp.StatusCode, body)
		return
	}

	log.Infof("notified admin: %s", buffer.String())
}
