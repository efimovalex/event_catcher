package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapEvent(t *testing.T) {

	testCases := []struct {
		testName string
		events   map[string]interface{}

		expectedResponse Event
	}{
		{
			testName: "Empty",
			events:   map[string]interface{}{},

			expectedResponse: Event{
				timeField: "timestamp",
			},
		},
		{
			testName: "All Fields - single category",
			events: map[string]interface{}{
				"status":           "true",
				"sg_event_id":      "sg_event_id",
				"sg_message_id":    "sg_message_id",
				"event":            "send",
				"email":            "email@email.com",
				"timestamp":        "1476629755",
				"smtp-id":          "id",
				"send_at":          "today",
				"reason":           "delay",
				"type":             "email",
				"tls":              "1",
				"cert_err":         "err",
				"ip":               "192.168.123.123",
				"url":              "http://url.com",
				"url_offset_index": 1,
				"url_offset_type":  "type",
				"asm_group_id":     12,
				"useragent":        "useragent",
				"ip_pool_name":     "name",
				"ip_pool_id":       221,
				"category":         "single_category",
				"newsletter": map[string]string{
					"newsletter_user_list_id": "11",
					"newsletter_send_id":      "123",
					"newsletter_id":           "321",
				},
				"marketing_campaign_id":  33,
				"nlvx_campaign_id":       111,
				"nlvx_campaign_split_id": 321,
				"nlvx_user_id":           432,
				"post_type":              "type_post",
				"unique_arg":             "arg",
			},

			expectedResponse: Event{
				timeField:        "timestamp",
				Status:           "true",
				SGEventID:        "sg_event_id",
				SGMessageID:      "sg_message_id",
				Event:            "send",
				Email:            "email@email.com",
				Timestamp:        1476629755,
				SMTPID:           "id",
				SendAt:           "today",
				Reason:           "delay",
				Type:             "email",
				TLS:              true,
				CertificateError: "err",
			},
		},
	}

	for _, test := range testCases {
		event := Event{}

		event.MapEvent(test.events)

		assert.Equal(t, test.expectedResponse, event)
	}
}
