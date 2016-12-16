// +build rest_acceptance

package rest_acceptance

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"github.com/sendgrid/reseller_integrations/helper/requester"
	"github.com/stretchr/testify/assert"
)

var DBAdaptor database.Adaptor

func TestConsumer(t *testing.T) {
	DBAdaptor := database.NewAdaptor(strings.Split(acceptanceConfig.CassandraInterfaces, ","), acceptanceConfig.CassandraUser, acceptanceConfig.CassandraPassword)

	restURL := fmt.Sprintf("http://%s:%d/v1/events", acceptanceConfig.Interface, acceptanceConfig.Port)

	req := requester.New(&http.Client{}, "consumer_tests")

	var response string

	resp, err := req.Make("GET", restURL, "", &response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Contains(t, string(resp.Body), `{"events":[]}`)

	event := database.Event{
		Event:          "processed",
		Email:          "john.doe@sendgrid.com",
		Timestamp:      time.Unix(1466519395, 0),
		SGEventID:      "2VzcPw12113x5Pv122137SdW2vUug4t-xKymw",
		SGMessageID:    "142d9f3f351.7618.254f56.filter-147.22649.52A663508.0",
		SMTPID:         "<142d9f3f351.7618.254f56@sendgrid.com>",
		Categories:     []string{"category1", "category2"},
		UniqueArgumets: map[string]string{"id": "1022", "purchase": "PO1452297845", "uid": "123456"},
	}

	DBAdaptor.Save(event)

	resp, err = req.Make("GET", restURL, "", &response)
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, string(resp.Body), `{"events":[{"sg_event_id":"2VzcPw12113x5Pv122137SdW2vUug4t-xKymw","sg_message_id":"142d9f3f351.7618.254f56.filter-147.22649.52A663508.0","event":"processed","email":"john.doe@sendgrid.com","timestamp":"2016-06-21T14:29:55Z","smtp_id":"\u003c142d9f3f351.7618.254f56@sendgrid.com\u003e","unique_arguments":{"id":"1022","purchase":"PO1452297845","uid":"123456"},"categories":["category1","category2"]}]}`)

	event = database.Event{
		Event:          "sent",
		Email:          "filtered@sendgrid.com",
		Timestamp:      time.Unix(1466519399, 0),
		SGEventID:      "2VzcPw12113x5Pv122137SdW2sdsd-xKymw",
		SGMessageID:    "142d9f3f351.2341.254f56.filter-147.22649.52A663508.0",
		SMTPID:         "<142d9f3f351.2355.254f56@sendgrid.com>",
		Categories:     []string{"category1"},
		UniqueArgumets: map[string]string{},
	}

	DBAdaptor.Save(event)

	restURL = restURL + "?field_name=email&field_value=filtered@sendgrid.com"

	resp, err = req.Make("GET", restURL, "", &response)
	if err != nil {
		t.Error(err)
	}

	assert.Contains(t, string(resp.Body), `{"events":[{"sg_event_id":"2VzcPw12113x5Pv122137SdW2sdsd-xKymw","sg_message_id":"142d9f3f351.2341.254f56.filter-147.22649.52A663508.0","event":"processed","email":"filtered@sendgrid.com","timestamp":"2016-06-21T14:29:55Z","smtp_id":"\u003c142d9f3f351.2355.254f56@sendgrid.com\u003e","categories":["category1"]}]}`)
}
