package database

import (
	"errors"
)

type Event interface {
	Save(*Adaptor) error

	MapEvent(eventMap map[string]interface{})
}

func (a *Adaptor) AddEvent(event map[string]interface{}) error {
	eventObj, err := getMappedEvent(event)
	if err != nil {
		return err
	}

	saveErr := eventObj.Save(a)

	return saveErr
}

func EventFactory(eventType string) (Event, error) {
	switch eventType {
	case "bounce":
		return &BounceEvent{}, nil
	case "deffered":
		return &DeferredEvent{}, nil
	case "delivered":
		return &DeliveredEvent{}, nil
	case "dropped":
		return &DroppedEvent{}, nil
	case "processed":
		return &ProcessedEvent{}, nil
	case "click":
		return &ClickEvent{}, nil
	case "open":
		return &OpenEvent{}, nil
	case "spamreport":
		return &SpamReportEvent{}, nil
	case "unsubscribe":
		return &UnsubscribeEvent{}, nil
	case "group_unsubscribe":
		return &GroupUnsubscribeEvent{}, nil
	case "group_resubscribe":
		return &GroupResubscribeEvent{}, nil
	default:
		return nil, errors.New("Undefined event: " + eventType)
	}
}

func getMappedEvent(event map[string]interface{}) (Event, error) {
	eventObj, err := EventFactory(event["event"].(string))
	if err != nil {
		return nil, err
	}
	eventObj.MapEvent(event)

	return eventObj, nil
}
