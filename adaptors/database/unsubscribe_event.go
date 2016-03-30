package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hailocab/gocassa"
)

type UnsubscribeEvent struct {
	SGMessageID string `cql:"sg_message_id"`
	Email       string `cql:"email"`
	Timestamp   int64  `cql:"timestamp"`
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

func (e *UnsubscribeEvent) MapEvent(eventMap map[string]interface{}) {
	e.SGMessageID = eventMap["sg_message_id"].(string)
	delete(eventMap, "sg_message_id")
	e.Event = eventMap["event"].(string)
	delete(eventMap, "event")
	e.Email = eventMap["email"].(string)
	delete(eventMap, "email")
	e.Timestamp, _ = strconv.ParseInt(string(eventMap["timestamp"].(json.Number)), 10, 64)
	delete(eventMap, "timestamp")

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

func (e *UnsubscribeEvent) Save(adaptor *Adaptor) error {
	unsubscribeTable := adaptor.Session.Table("unsubscribe_events", &UnsubscribeEvent{}, gocassa.Keys{
		PartitionKeys: []string{"sg_event_id"},
	})

	if err := unsubscribeTable.Set(e).Run(); err != nil {
		return fmt.Errorf("Save error: %s", err.Error())
	}

	return nil
}
