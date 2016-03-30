package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hailocab/gocassa"
)

type OpenEvent struct {
	Email       string `cql:"email"`
	Timestamp   int64  `cql:"timestamp"`
	IP          string `cql:"ip"`
	SGEventID   string `cql:"sg_event_id"`
	SGMessageID string `cql:"sg_message_id"`
	UserAgent   string `cql:"user_agent"`
	Event       string `cql:"event"`

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

func (e *OpenEvent) MapEvent(eventMap map[string]interface{}) {
	e.SGEventID = eventMap["sg_event_id"].(string)
	delete(eventMap, "sg_event_id")
	e.SGMessageID = eventMap["sg_message_id"].(string)
	delete(eventMap, "sg_message_id")
	e.Event = eventMap["event"].(string)
	delete(eventMap, "event")
	e.Timestamp, _ = strconv.ParseInt(string(eventMap["timestamp"].(json.Number)), 10, 64)
	delete(eventMap, "timestamp")
	e.IP = eventMap["ip"].(string)
	delete(eventMap, "ip")
	e.UserAgent = eventMap["user_agent"].(string)
	delete(eventMap, "user_agent")

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

func (e *OpenEvent) Save(adaptor *Adaptor) error {
	openTable := adaptor.Session.Table("bounce_events", &OpenEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := openTable.Set(e).Run(); err != nil {
		return fmt.Errorf("Save error: %s", err.Error())
	}

	return nil
}
