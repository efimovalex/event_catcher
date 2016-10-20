package cache

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/redis.v3"

	"github.com/efimovalex/EventKitAPI/adaptors/database"
	"github.com/efimovalex/EventKitAPI/common"
)

// ServiceInterface is an abstraction of redis commands
type ServiceInterface interface {
	GetSet(string, []byte, int) ([]byte, error)
	Del(string) error
	SaveEventRequest([]byte, common.ListResponse, time.Duration)
	GetEventRequest([]byte) string
}

// Adaptor for redis abstraction layer
type Adaptor struct {
	client *redis.Client
	ttl    time.Duration
}

// New instantiates the redis abstraction layer
func NewAdaptor(addrs, password string) *Adaptor {
	client := redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return &Adaptor{client: client, ttl: 10 * time.Minute}
}

// Get returns original value
// See http://redis.io/commands/get
func (a *Adaptor) Get(key string) (string, error) {
	dataResponse, err := a.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return dataResponse, nil
}

// Set sets new value, and also sets ttl
// See http://redis.io/commands/set
func (a *Adaptor) Set(key string, value []byte, ttl time.Duration) error {
	_, err := a.client.Set(key, value, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the key/value
// See http://redis.io/commands/del
func (a *Adaptor) Delete(keys ...string) error {
	_, err := a.client.Del(keys...).Result()

	return err
}

// Incr increments the version of the key
// See http://redis.io/commands/incr
func (a *Adaptor) Incr(key string) (int64, error) {
	return a.client.Incr(key).Result()
}

// SaveEventRequest saves the Response struct in Redis.
func (a *Adaptor) SaveEventRequest(queryString []byte, response common.ListResponse) {
	requestHash := sha1.Sum(queryString)

	version, err := a.Incr(string(requestHash[:]))
	if err != nil {
		log.Println("error incrementing request version: " + err.Error())
	}

	versionedKey := fmt.Sprintf("%s:%d", requestHash, version)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("error encoding response: " + err.Error())
	}

	if err := a.Set(versionedKey, jsonResponse, a.ttl); err != nil {
		log.Println("error saving response in cache: " + err.Error())
	}

	go a.saveEvents(response.Events)
}

func (a *Adaptor) saveEvents(events []database.Event) {
	for _, event := range events {
		jsonResponse, err := json.Marshal(event)
		if err != nil {
			log.Println("error encoding response: " + err.Error())
		}
		if err := a.Set(event.SGEventID, jsonResponse, a.ttl); err != nil {
			log.Println("error saving response in cache: " + err.Error())
		}
	}
}

// GetEventRequest retrieves the JSON string Response saved in cache
func (a *Adaptor) GetEventRequest(queryString []byte) string {
	requestHash := sha1.Sum(queryString)
	stringRequestHash := string(requestHash[:])
	currentVersion, err := a.Get(stringRequestHash)
	if err != nil {
		log.Println("error getting response version from cache: " + err.Error())
	}

	versionedKey := fmt.Sprintf("%s:%s", stringRequestHash, currentVersion)

	log.Println(versionedKey)
	response, err := a.Get(versionedKey)
	if err != nil {
		log.Println("error getting response from cache: " + err.Error())
	}

	return response
}

// GetEvent retrieves the JSON string event saved in cache
func (a *Adaptor) GetEvent(eventID string) string {
	event, err := a.Get(eventID)
	if err != nil {
		log.Println("error getting response from cache: " + err.Error())
	}

	return event
}
