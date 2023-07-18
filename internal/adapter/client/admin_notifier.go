package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var notificationServiceURL = os.Getenv("GREENBONE_NOTIFICATION_URL")

type Alarm struct {
	Level   string `json:"level"`
	Abbr    string `json:"employeeAbbreviation"`
	Message string `json:"message"`
}

// NotifyAdmin sends an HTTP request to an external notification service in order to inform
// an admin about employees with excessive computer demands
func NotifyAdmin(employeeAbbr string, computerCount int64) {
	alarm := Alarm{
		Level:   "warning",
		Abbr:    employeeAbbr,
		Message: fmt.Sprintf("employee %s has %d computers", employeeAbbr, computerCount),
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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("failed to close response body from admin notification service: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body from admin notification service: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("request to admin notification service failed (%d): %q", resp.StatusCode, body)
		return
	}

	log.Infof("notified admin: %s", body)
}
