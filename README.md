# EventKitAPI v2

## Event Consumer

### Starting the service 

```
$ ./bin/start_consumer

### Consumer Endpoint

#### POST /v1/events

```
curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '[
    {
        "email": "john.doe+3@sendgrid.com",
        "sg_event_id": "2VzcPw12113x5Pv2137SdW2vUug4t-xKymw",
        "sg_message_id": "142d9f3f351.7618.254f56.filter-147.22649.52A663508.0",
        "timestamp": 1466519395,
        "smtp-id": "<142d9f3f351.7618.254f56@sendgrid.com>",
        "event": "processed",
        "category":["category1", "category2"],
        "id": "1022",
        "purchase":"PO1452297845",
        "uid": "123456"
    }
]' "http://192.168.100.2:59000/v1/events"
```