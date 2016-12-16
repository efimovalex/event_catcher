# EventKitAPI v2

## Required External Services:

### Cassandra database: 

Edit example.env and replace the following settings with the Cassandra connection info

```bash
# Consumer API settings
export CONSUMER_CASSANDRA_INTERFACES="0.0.0.0"
export CONSUMER_CASSANDRA_USERNAME=""
export CONSUMER_CASSANDRA_PASSWORD=""

# REST api settings
export REST_CASSANDRA_INTERFACES="0.0.0.0"
export REST_CASSANDRA_USERNAME=""
export REST_CASSANDRA_PASSWORD=""
```

### Redis [Optional]

To enable Redis caching of REST responses set the following ENV variable:
```bash
export REST_ENABLE_CACHING=1
```
Then edit the settings for Redis:
```bash
export REST_CACHE_REDIS_URL="0.0.0.0:6379"
export REST_CACHE_REDIS_USERNAME=""
export REST_CACHE_REDIS_PASSWORD=""
```

## Event Consumer API

### Service description
The scope of this application is to capture events from the SendGrid Event WebHook. 
SendGrid sends events in bulk when the number of events on a given server for a user reaches 100.

### Building the service
```
$ ./bin/build
```

### Starting the service 
example.env and replace the following settings 

```
export CONSUMER_INTERFACE="0.0.0.0"
export CONSUMER_PORT=59000
export CONSUMER_MAX_JOB_QUEUE="20"
export CONSUMER_MAX_WORKER="20"
```

```
$ source example.env
$ ./build/event_kit consumer
```

### Consumer Endpoint

#### POST /v1/events

```
curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '
[
    {
        "email": "john.doe+1@sendgrid.com",
        "sg_event_id": "2VzcPw12113x5Pv2137SdW2vUug4t-xKymw",
        "sg_message_id": "f3f351142d9.7618.254f56.filter-147.22649.52A663508.0",
        "timestamp": 1466519395,
        "smtp-id": "<f3f351142d9.7618.254f56@sendgrid.com>",
        "event": "processed",
        "category":["category1", "category2"],
        "id": "11022",
        "purchase":"PO1452297845"
    },
    {
        "email": "john.doe+2@sendgrid.com",
        "sg_event_id": "sd344w12113x5Pv2137SdW2vUug4t-xKymw",
        "sg_message_id": "254f56.7618.254f56.filter-147.22649.52A663508.0",
        "timestamp": 1466519395,
        "smtp-id": "<254f56.7618.254f56@sendgrid.com>",
        "event": "processed",
        "category":["category1", "category2"],
        "id": "12022",
        "purchase":"PO1452297845"
    },
    {
        "email": "john.doe+3@sendgrid.com",
        "sg_event_id": "dhfh43242d9f314f351SdW2vUug4t-xKymw",
        "sg_message_id": "2d9f314f351.7618.254f56.filter-147.22649.52A663508.0",
        "timestamp": 1466519395,
        "smtp-id": "<2d9f314f351.7618.254f56@sendgrid.com>",
        "event": "processed",
        "category":["1", "d3"],
        "id": "10222",
        "purchase":"PO1452297845"
    }
]' "http://localhost:59000/v1/events"
```

##### Responses 

```
HTTP/1.1 202 Accepted
Date: Fri, 16 Dec 2016 10:19:36 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

## Event REST API

### Service description
The scope of this application is to expose trough sa set of endpoints CRUD operations for the events captured with the Event Consumer API that are indexed in Cassandra.

### Building the service
```
$ ./bin/build
```

### Starting the service 
example.env and replace the following settings 

```
export CONSUMER_INTERFACE="0.0.0.0"
export CONSUMER_PORT=59000
export CONSUMER_MAX_JOB_QUEUE="20"
export CONSUMER_MAX_WORKER="20"
```

```
$ source example.env
$ ./build/event_kit consumer
```

### Consumer Endpoint

#### GET /v1/events

```
curl -X GET "http://localhost:59001/v1/events?field_name=email&field_value=john.doe+1@sendgrid.com&limit=4&offset_id=VzcPwxPv7SdWvUugt-xKymw2423%22"
```
```
curl -X GET "http://localhost:59001/v1/events?field_name=day&start_date=2013-12-09&end_date=2013-12-09"
```

Parameter | Value | Required
--- | --- | ---
field_name | email/sg_message_id/event/day | yes
field_value | Value of the selected type above | no
start_date | Date value if field_name = day | no
end_date| Date value if field_name = day | no
limit | Number of results returned | no
offset_id | Id of the last event returned in the previous page | no 

##### Responses 
```
HTTP/1.1 200 OK
Cache-Control: private, max-age=600
Content-Type: application/json; charset=utf-8
Vary: Accept-Encoding
Date: Fri, 16 Dec 2016 10:21:03 GMT
Content-Length: 13

{"events":
    [
        {
            "sg_event_id": "2VzcPw12113x5Pv2137SdW2vUug4t-xKymw",
            "sg_message_id": "f3f351142d9.7618.254f56.filter-147.22649.52A663508.0",
            "event": "processed",
            "email": "john.doe+1@sendgrid.com",
            "timestamp": "2016-06-21T14:29:55Z",
            "smtp_id": "<f3f351142d9.7618.254f56@sendgrid.com>",
            "unique_arguments": {
                "id": "11022",
                "purchase": "PO1452297845"
            },
            "categories": [
                "category1",
                "category2"
            ]
        }
    ]
}
```

#### GET /v1/event/:sg_event_id

```
curl -X GET "http://localhost:59001/v1/event/sd344w12113x5Pv2137SdW2vUug4t-xKymw
```

##### Responses 
```
HTTP/1.1 200 OK
Cache-Control: private, max-age=600
Content-Type: application/json; charset=utf-8
Vary: Accept-Encoding
Date: Fri, 16 Dec 2016 10:21:03 GMT
Content-Length: 13

{
    "sg_event_id": "2VzcPw12113x5Pv2137SdW2vUug4t-xKymw",
    "sg_message_id": "f3f351142d9.7618.254f56.filter-147.22649.52A663508.0",
    "event": "processed",
    "email": "john.doe+1@sendgrid.com",
    "timestamp": "2016-06-21T14:29:55Z",
    "smtp_id": "<f3f351142d9.7618.254f56@sendgrid.com>",
    "unique_arguments": {
        "id": "11022",
        "purchase": "PO1452297845"
    },
    "categories": [
        "category1",
        "category2"
    ]
}
```

#### DELETE /v1/event/:sg_event_id

```
curl -X DELETE "http://localhost:59001/v1/event/sd344w12113x5Pv2137SdW2vUug4t-xKymw
```

##### Responses 
```
HTTP/1.1 200 OK
Cache-Control: private, max-age=600
Content-Type: application/json; charset=utf-8
Vary: Accept-Encoding
Date: Fri, 16 Dec 2016 10:21:03 GMT
Content-Length: 13
```
