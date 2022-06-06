# HTTP-Multiplexer

#### Multiplexer is a tiny service that  handles HTTP requests. Service waits for request consisting of some number of URLs, make concurrent HTTP request calls and sends received data back to the Caller.

# API Description
## 1. Get info from provided URls
#### The one and only endpoint 
#### Query: ```POST ServiceURL/multiply ```
#### Body: 
```
{
    "urls": [
    "some number of URLs",
    "url1",
    "url2",
    "url3",
    ...
    ]
}
```
#### Response on success:
```
{
    "data": [
        {
            "url": "some number of URLs",
            "response": "json response from url encoded in Base64"
        },
        {
            "url": "url1",
            "response": "url1 json response from GET request in Base64"
        },
        {
            "url": "url2",
            "response": "url2 json response from GET request in Base64"
        }
    ]
}
```
#### Response on error:
#### Case 1: Bad request provided:
Consists of ```Bad Request``` Header and empty JSON body

#### Case 2: An error occurred:
 - Any call from provided URLs timeouted
```
{
    "error": "url some_heavy_endpoint: request takes too long"
}
```
- Any call from provided URLs errored
```
{
    "error": "url some_heavy_endpoint: err_msg_from_requested_URL"
}
```
# Service features
 - Service handles requests consisting of 20 URLs at max. Otherwise it will return a corresponding ErrMessage in JSON.
 - Response from provided URLs from request is encoded using Base64.
 - Service handles 100 concurrent incoming requests at max. The rest of requests will wait in queue.
 - Service handles 4 concurrent outgoing requests per provided URL at max. The rest of requests will wait in queue.
 - Service timeout on each outgoing URL request is set to 1 second
 - Service provides Graceful Shutdown. That is, when service gets `SIGTERM` singal for example, it will not receive any new request, but will handle ones that are currently registered in queue, and shutdown only after.   
[Soon] You will be able to change this variables in `values.yaml` config file
# To launch service locally:
1. git clone this repo
```
git clone https://github.com/hyonosake/HTTP-Multiplexer && cd HTTP-Multiplexer/
```
2. Use docker-compose to build containers with PostgreSQL and API service
```
docker-compose up
```
3. Open your browser or API-testing app, enter ```localhost:1234/muliply``` and make some requests!
## [WIP] Testing
Static tests are yet to be done. For now, you can run service and test it using curl or Postman, for example

Some Test cases:

Request
```
{
    "urls": [
        "http://google.com",
        "http://google.com",
        "http://google.com",
        "http://google.com",
        "http://google.com",
        "http://google.com"
    ]
}
```
Response:
```
{
    "data": [
        {
            "url": "http://google.com",
            "response": "some_data"
        },
        {
            "url": "http://google.com",
            "response": "some_data"
        },
        {
            "url": "http://google.com",
            "response": "some_data"
        },
        {
            "url": "http://google.com",
            "response": "some_data"
        },
        {
            "url": "http://google.com",
            "response": "some_data"
        },
        {
            "url": "http://google.com",
            "response": "some_data"
        },
    ]
}
```

Request
```
{
    "urls": [
        "http://google.com",
        "http://yahoo.com"
    ]
}
```
Response:
```
{
    "error": "url http://yahoo.com: request takes too long"
}
```







