// +build consumer_acceptance

package consumer_acceptance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sendgrid/reseller_integrations/helper/requester"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	consumerURL := fmt.Sprintf("http://%s:%d/v1/events", acceptanceConfig.Interface, acceptanceConfig.Port)
	body := `
[
    {
        "email": "john.doe@sendgrid.com",
        "sg_event_id": "2VzcPw12113x5Pv122137SdW2vUug4t-xKymw",
        "sg_message_id": "142d9f3f351.7618.254f56.filter-147.22649.52A663508.0",
        "timestamp": 1466519395,
        "smtp-id": "<142d9f3f351.7618.254f56@sendgrid.com>",
        "event": "processed",
        "category":["category1", "category2"],
        "id": "1022",
        "purchase":"PO1452297845",
        "uid": "123456"
    }
]`

	req := requester.New(&http.Client{}, "consumer_tests")

	var response string

	// Test Create NLvX addon call
	resp, err := req.Make("POST", consumerURL, body, &response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusAccepted, resp.Status)
}
