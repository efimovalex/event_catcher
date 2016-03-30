package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hailocab/gocassa"
)

type BounceEvent struct {
	Status           string `cql:"status"`
	SGEventID        string `cql:"sg_event_id"`
	SGMessageID      string `cql:"sg_message_id"`
	Event            string `cql:"event"`
	Email            string `cql:"email"`
	Timestamp        int64  `cql:"timestamp"`
	SMTPID           string `cql:"smtp_id"`
	Reason           string `cql:"reason"`
	Type             string `cql:"type"`
	IP               string `cql:"ip"`
	TLS              bool   `cql:"tls"`
	CertificateError string `cql:"cert_err"`

	// optional fields
	NewsletterUserListID int `cql:"newsletter_user_list_id"`
	NewsletterID         int `cql:"newsletter_id"`
	NewsletterSendID     int `cql:"newsletter_send_id"`

	MarketingCampainName    string `cql:"marketing_campain_name"`
	MarketingCampainID      int    `cql:"marketing_campain_id"`
	MarketingCampainVersion string `cql:"marketing_campain_version"`
	NLVXCampainID           int    `cql:"nlvx_campain_id"`
	NLVXUserID              int    `cql:"nlvx_user_id"`
	NLVXCampainSplitID      int    `cql:"nlvx_campain_split_id"`
	PostType                string `cql:"post_type"`

	ASMGroupID     int               `cql:"asm_group_id"`
	UniqueArgumets map[string]string `cql:"unique_arguments"`
	Categories     []string          `cql:"categories"`
}

func (e *BounceEvent) MapEvent(eventMap map[string]interface{}) {
	e.Status = eventMap["status"].(string)
	delete(eventMap, "status")
	e.SGEventID = eventMap["sg_event_id"].(string)
	delete(eventMap, "sg_event_id")
	e.SGMessageID = eventMap["sg_message_id"].(string)
	delete(eventMap, "sg_message_id")
	e.Event = eventMap["event"].(string)
	delete(eventMap, "event")
	e.Email = eventMap["email"].(string)
	delete(eventMap, "email")
	e.Timestamp, _ = strconv.ParseInt(string(eventMap["timestamp"].(json.Number)), 10, 64)
	delete(eventMap, "timestamp")
	e.SMTPID = eventMap["smtp-id"].(string)
	delete(eventMap, "smtp-id")
	e.Reason = eventMap["reason"].(string)
	delete(eventMap, "reason")
	e.Type = eventMap["type"].(string)
	delete(eventMap, "type")
	e.TLS = (eventMap["tls"] == "1")
	delete(eventMap, "tls")

	if val, ok := eventMap["cert_err"]; ok {
		e.CertificateError = val.(string)
		delete(eventMap, "cert_err")
	}

	if val, ok := eventMap["ip"]; ok {
		e.IP = val.(string)
		delete(eventMap, "ip")
	}

	if val, ok := eventMap["asm_group_id"]; ok {
		e.ASMGroupID = val.(int)
		delete(eventMap, "asm_group_id")
	}

	if _, ok := eventMap["category"]; ok {
		switch reflect.TypeOf(eventMap["category"]).Kind() {
		case reflect.Slice:
			e.Categories = eventMap["category"].([]string)
		case reflect.String:
			e.Categories = []string{eventMap["category"].(string)}
		}
		delete(eventMap, "category")
	}

	if val, ok := eventMap["newsletter"]; ok {
		nl := val.(map[string]string)

		e.NewsletterUserListID, _ = strconv.Atoi(nl["newsletter_user_list_id"])
		e.NewsletterSendID, _ = strconv.Atoi(nl["newsletter_send_id"])
		e.NewsletterID, _ = strconv.Atoi(nl["newsletter_id"])

		delete(eventMap, "newsletter")
	}

	if val, ok := eventMap["marketing_campaign_id"]; ok {
		e.MarketingCampainID = val.(int)
		delete(eventMap, "marketing_campaign_id")
	}
	if val, ok := eventMap["nlvx_campaign_id"]; ok {
		e.NLVXCampainID = val.(int)
		delete(eventMap, "nlvx_campaign_id")
	}
	if val, ok := eventMap["nlvx_campaign_split_id"]; ok {
		e.NLVXCampainSplitID = val.(int)
		delete(eventMap, "nlvx_campaign_split_id")
	}
	if val, ok := eventMap["nlvx_user_id"]; ok {
		e.NLVXUserID = val.(int)
		delete(eventMap, "nlvx_user_id")
	}
	if val, ok := eventMap["post_type"]; ok {
		e.PostType = val.(string)
		delete(eventMap, "post_type")
	}

	for index, val := range eventMap {
		e.UniqueArgumets[index] = val.(string)
	}
}

func (e *BounceEvent) Save(adaptor *Adaptor) error {
	bounceTable := adaptor.Session.Table("bounce_events", &BounceEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := bounceTable.Set(e).Run(); err != nil {
		return fmt.Errorf("Save error: %s", err.Error())
	}

	return nil
}
