package database

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

type Event struct {
	timeField string
	Status    string `cql:"status,omitempty" json:"status,omitempty"`

	SGEventID   string    `cql:"sg_event_id" json:"sg_event_id,omitempty"`
	SGMessageID string    `cql:"sg_message_id" json:"sg_message_id,omitempty"`
	Event       string    `cql:"event" json:"event,omitempty"`
	Email       string    `cql:"email" json:"email,omitempty"`
	Timestamp   time.Time `cql:"timestamp" json:"timestamp,omitempty"`
	SMTPID      string    `cql:"smtp_id" json:"smtp_id,omitempty"`

	SendAt string `cql:"send_at" json:"send_at,omitempty"`

	Reason string `cql:"reason" json:"reason,omitempty"`
	Type   string `cql:"type" json:"type,omitempty"`

	IP string `cql:"ip" json:"ip,omitempty"`

	TLS              bool   `cql:"tls" json:"tls,omitempty"`
	CertificateError string `cql:"cert_err" json:"cert_err,omitempty"`

	Response string `cql:"response" json:"response,omitempty"`
	Attempt  string `cql:"attempt" json:"attempt,omitempty"`

	URL            string `cql:"url" json:"url,omitempty"`
	URLOffsetIndex int    `cql:"url_offset_index" json:"url_offset_index,omitempty"`
	URLOffsetType  string `cql:"url_offset_type" json:"url_offset_type,omitempty"`

	UserAgent string `cql:"user_agent" json:"user_agent,omitempty"`

	IPPoolName string `cql:"ip_pool_name" json:"ip_pool_name,omitempty"`
	IPPoolID   int    `cql:"ip_pool_id" json:"ip_pool_id,omitempty"`

	NewsletterUserListID int `cql:"newsletter_user_list_id,omitempty" json:"newsletter_user_list_id,omitempty"`
	NewsletterID         int `cql:"newsletter_id" json:"newsletter_id,omitempty"`
	NewsletterSendID     int `cql:"newsletter_send_id" json:"newsletter_send_id,omitempty"`

	MarketingCampainName    string `cql:"marketing_campain_name,omitempty" json:"marketing_campain_name,omitempty"`
	MarketingCampainID      int    `cql:"marketing_campain_id,omitempty" json:"marketing_campain_id,omitempty"`
	MarketingCampainVersion string `cql:"marketing_campain_version,omitempty" json:"marketing_campain_version,omitempty"`
	NLVXCampainID           int    `cql:"nlvx_campain_id,omitempty" json:"nlvx_campain_id,omitempty"`
	NLVXUserID              int    `cql:"nlvx_user_id,omitempty" json:"nlvx_user_id,omitempty"`
	NLVXCampainSplitID      int    `cql:"nlvx_campain_split_id,omitempty" json:"nlvx_campain_split_id,omitempty"`
	PostType                string `cql:"post_type,omitempty" json:"post_type,omitempty"`

	ASMGroupID     int               `cql:"asm_group_id,omitempty" json:"asm_group_id,omitempty"`
	UniqueArgumets map[string]string `cql:"unique_arguments,omitempty" json:"unique_arguments,omitempty"`
	Categories     []string          `cql:"categories,omitempty" json:"categories,omitempty"`
}

func (e *Event) MapEvent(eventMap map[string]interface{}) {
	e.timeField = "timestamp"
	if val, ok := eventMap["status"]; ok {
		e.Status = val.(string)
		delete(eventMap, "status")
	}

	if val, ok := eventMap["sg_event_id"]; ok {
		e.SGEventID = val.(string)
		delete(eventMap, "sg_event_id")
	}

	if val, ok := eventMap["sg_message_id"]; ok {
		e.SGMessageID = val.(string)
		delete(eventMap, "sg_message_id")
	}

	if val, ok := eventMap["event"]; ok {
		e.Event = val.(string)
		delete(eventMap, "event")
	}

	if val, ok := eventMap["email"]; ok {
		e.Email = val.(string)
		delete(eventMap, "email")
	}

	if val, ok := eventMap["timestamp"]; ok {
		intTime, _ := strconv.ParseInt(string(val.(json.Number)), 10, 64)
		e.Timestamp = time.Unix(intTime, 0)
		delete(eventMap, "timestamp")
	}

	if val, ok := eventMap["smtp-id"]; ok {
		e.SMTPID = val.(string)
		delete(eventMap, "smtp-id")
	}

	if val, ok := eventMap["send_at"]; ok {
		e.SendAt = val.(string)
		delete(eventMap, "send_at")
	}

	if val, ok := eventMap["reason"]; ok {
		e.Reason = val.(string)
		delete(eventMap, "reason")
	}

	if val, ok := eventMap["type"]; ok {
		e.Type = val.(string)
		delete(eventMap, "type")
	}

	if val, ok := eventMap["tls"]; ok {
		e.TLS = (val == "1")
		delete(eventMap, "tls")
	}

	if val, ok := eventMap["cert_err"]; ok {
		e.CertificateError = val.(string)
		delete(eventMap, "cert_err")
	}

	if val, ok := eventMap["ip"]; ok {
		e.IP = val.(string)
		delete(eventMap, "ip")
	}

	if val, ok := eventMap["url"]; ok {
		e.URL = val.(string)
		delete(eventMap, "url")
	}

	if val, ok := eventMap["url_offset_index"]; ok {
		e.URLOffsetIndex = val.(int)
		delete(eventMap, "url_offset_index")
	}

	if val, ok := eventMap["url_offset_type"]; ok {
		e.URLOffsetType = val.(string)
		delete(eventMap, "url_offset_type")
	}

	if val, ok := eventMap["asm_group_id"]; ok {
		e.ASMGroupID = val.(int)
		delete(eventMap, "asm_group_id")
	}

	if val, ok := eventMap["useragent"]; ok {
		e.UserAgent = val.(string)
		delete(eventMap, "useragent")
	}

	if val, ok := eventMap["ip_pool_name"]; ok {
		e.IPPoolName = val.(string)
		delete(eventMap, "ip_pool_name")
	}

	if val, ok := eventMap["ip_pool_id"]; ok {
		e.IPPoolID = val.(int)
		delete(eventMap, "ip_pool_id")
	}

	if _, ok := eventMap["category"]; ok {
		switch reflect.TypeOf(eventMap["category"]).Kind() {
		case reflect.Slice:
			categories := interfaceSlice(eventMap["category"])
			for _, val := range categories {
				e.Categories = append(e.Categories, val.(string))
			}
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

	e.UniqueArgumets = make(map[string]string)

	for index, val := range eventMap {
		e.UniqueArgumets[index] = val.(string)
	}

	if len(e.UniqueArgumets) == 0 {
		e.UniqueArgumets = nil
	}

	if len(e.Categories) == 0 {
		e.Categories = nil
	}
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
