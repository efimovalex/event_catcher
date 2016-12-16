package database

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAdaptor(t *testing.T) {
	adaptor := NewAdaptor([]string{"localhost"}, "", "")

	assert.Equal(t, adaptor.urls, []string{"localhost"})
	assert.Equal(t, adaptor.eventTimeDateRange, 7*24*time.Hour)
	assert.Equal(t, adaptor.Username, "")
	assert.Equal(t, adaptor.Password, "")
}

func TestAddEvent(t *testing.T) {
	adaptor := NewAdaptor([]string{"localhost"}, "", "")

	eventMap := map[string]interface{}{}

	result := adaptor.AddEvent(eventMap)

	assert.Equal(t, result.Error(), "Save error: Key may not be empty")

	eventMap = map[string]interface{}{
		"email":         "john.doe@sendgrid.com",
		"sg_event_id":   "2VzcPw12113x5Pv122137SdW2vUug4t-xKymw",
		"sg_message_id": "142d9f3f351.7618.254f56.filter-147.22649.52A663508.0",
		"timestamp":     json.Number("1466519395"),
		"smtp-id":       "<142d9f3f351.7618.254f56@sendgrid.com>",
		"event":         "processed",
		"category":      []string{"category1", "category2"},
		"id":            "1022",
		"purchase":      "PO1452297845",
		"uid":           "123456",
	}

	result = adaptor.AddEvent(eventMap)

	assert.Equal(t, result, nil)

}
